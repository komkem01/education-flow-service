package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/app/utils/hashing"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

const defaultAccessTokenTTL = 24 * time.Hour

type Service struct {
	tracer         trace.Tracer
	db             serviceDB
	appKey         []byte
	accessTokenTTL time.Duration
}

type serviceDB interface {
	entitiesinf.MemberEntity
	entitiesinf.MemberRoleEntity
}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     serviceDB
	appKey string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken string
	ExpiresAt   time.Time
	Member      *ent.Member
	Roles       []ent.MemberRole
}

type SwitchRoleInput struct {
	Claims *TokenClaims
	Role   ent.MemberRole
}

type SwitchRoleResult struct {
	AccessToken string
	ExpiresAt   time.Time
	Role        ent.MemberRole
	Roles       []ent.MemberRole
}

type TokenClaims struct {
	MemberID  uuid.UUID
	SchoolID  uuid.UUID
	Role      ent.MemberRole
	Roles     []ent.MemberRole
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type tokenPayload struct {
	Sub   string   `json:"sub"`
	Sid   string   `json:"sid"`
	Role  string   `json:"role"`
	Roles []string `json:"roles,omitempty"`
	Iat   int64    `json:"iat"`
	Exp   int64    `json:"exp"`
}

var (
	ErrInvalidCredentials = errors.New("invalid-credentials")
	ErrInactiveMember     = errors.New("inactive-member")
	ErrInvalidToken       = errors.New("invalid-token")
	ErrExpiredToken       = errors.New("expired-token")
	ErrRoleNotAllowed     = errors.New("role-not-allowed")
)

func newService(opt *Options) *Service {
	ttl := defaultAccessTokenTTL
	if opt.Config != nil && opt.Config.Val != nil && opt.Config.Val.AccessTokenTTLMinutes > 0 {
		ttl = time.Duration(opt.Config.Val.AccessTokenTTLMinutes) * time.Minute
	}

	key := strings.TrimSpace(opt.appKey)
	if key == "" {
		key = "education-flow-default-secret"
	}

	return &Service{
		tracer:         opt.tracer,
		db:             opt.db,
		appKey:         []byte(key),
		accessTokenTTL: ttl,
	}
}

func (s *Service) Login(ctx context.Context, input *LoginInput) (*LoginResult, error) {
	member, err := s.db.GetMemberByEmail(ctx, strings.TrimSpace(strings.ToLower(input.Email)))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !member.IsActive {
		return nil, ErrInactiveMember
	}

	if !hashing.CheckPasswordHash([]byte(member.Password), []byte(strings.TrimSpace(input.Password))) {
		return nil, ErrInvalidCredentials
	}

	now := time.Now().UTC()
	expiresAt := now.Add(s.accessTokenTTL)
	roles, err := s.db.ListMemberRolesByMemberID(ctx, member.ID)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, ErrInvalidCredentials
	}

	member.Role = roles[0]

	token, err := s.generateAccessToken(member, roles, now, expiresAt)
	if err != nil {
		return nil, err
	}

	if err := s.db.UpdateMemberLastLoginByID(ctx, member.ID); err != nil {
		return nil, err
	}

	member.LastLogin = &now

	return &LoginResult{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		Member:      member,
		Roles:       roles,
	}, nil
}

func (s *Service) ParseAccessToken(token string) (*TokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	headerPayload := parts[0] + "." + parts[1]
	expectedSig := s.sign(headerPayload)
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return nil, ErrInvalidToken
	}

	payloadRaw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	var payload tokenPayload
	if err := json.Unmarshal(payloadRaw, &payload); err != nil {
		return nil, ErrInvalidToken
	}

	memberID, err := uuid.Parse(payload.Sub)
	if err != nil {
		return nil, ErrInvalidToken
	}
	schoolID, err := uuid.Parse(payload.Sid)
	if err != nil {
		return nil, ErrInvalidToken
	}

	expiresAt := time.Unix(payload.Exp, 0).UTC()
	if time.Now().UTC().After(expiresAt) {
		return nil, ErrExpiredToken
	}

	roles := make([]ent.MemberRole, 0, len(payload.Roles)+1)
	for _, role := range payload.Roles {
		parsed, ok := parseKnownMemberRole(role)
		if !ok {
			return nil, ErrInvalidToken
		}
		exists := false
		for _, existing := range roles {
			if existing == parsed {
				exists = true
				break
			}
		}
		if !exists {
			roles = append(roles, parsed)
		}
	}
	if payload.Role != "" {
		parsedPrimary, ok := parseKnownMemberRole(payload.Role)
		if !ok {
			return nil, ErrInvalidToken
		}
		exists := false
		for _, existing := range roles {
			if existing == parsedPrimary {
				exists = true
				break
			}
		}
		if !exists {
			roles = append([]ent.MemberRole{parsedPrimary}, roles...)
		}
	}

	if len(roles) == 0 {
		return nil, ErrInvalidToken
	}

	return &TokenClaims{
		MemberID:  memberID,
		SchoolID:  schoolID,
		Role:      roles[0],
		Roles:     roles,
		IssuedAt:  time.Unix(payload.Iat, 0).UTC(),
		ExpiresAt: expiresAt,
	}, nil
}

func (s *Service) SwitchRole(ctx context.Context, input *SwitchRoleInput) (*SwitchRoleResult, error) {
	if input == nil || input.Claims == nil {
		return nil, ErrInvalidToken
	}

	memberRole, ok := parseKnownMemberRole(string(input.Role))
	if !ok {
		return nil, ErrRoleNotAllowed
	}

	member, err := s.db.GetMemberByID(ctx, input.Claims.MemberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	if member.SchoolID != input.Claims.SchoolID || !member.IsActive {
		return nil, ErrInvalidToken
	}

	roles, err := s.db.ListMemberRolesByMemberID(ctx, member.ID)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, ErrInvalidToken
	}

	validRoles := normalizeKnownRoles(roles)
	if !containsRole(validRoles, memberRole) {
		return nil, ErrRoleNotAllowed
	}

	orderedRoles := orderRolesWithPrimary(validRoles, memberRole)
	now := time.Now().UTC()
	expiresAt := now.Add(s.accessTokenTTL)
	member.Role = memberRole

	accessToken, err := s.generateAccessToken(member, orderedRoles, now, expiresAt)
	if err != nil {
		return nil, err
	}

	return &SwitchRoleResult{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
		Role:        memberRole,
		Roles:       orderedRoles,
	}, nil
}

func (s *Service) generateAccessToken(member *ent.Member, roles []ent.MemberRole, now, expiresAt time.Time) (string, error) {
	headerBytes, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", err
	}

	roleStrings := make([]string, 0, len(roles))
	for _, role := range roles {
		roleStrings = append(roleStrings, string(role))
	}

	payloadBytes, err := json.Marshal(tokenPayload{
		Sub:   member.ID.String(),
		Sid:   member.SchoolID.String(),
		Role:  string(primaryRole(roles)),
		Roles: roleStrings,
		Iat:   now.Unix(),
		Exp:   expiresAt.Unix(),
	})
	if err != nil {
		return "", err
	}

	head := base64.RawURLEncoding.EncodeToString(headerBytes)
	body := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signed := head + "." + body
	signature := s.sign(signed)

	return fmt.Sprintf("%s.%s", signed, signature), nil
}

func (s *Service) sign(data string) string {
	h := hmac.New(sha256.New, s.appKey)
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func primaryRole(roles []ent.MemberRole) ent.MemberRole {
	if len(roles) > 0 {
		return roles[0]
	}

	return ent.MemberRoleAdmin
}

func parseKnownMemberRole(value string) (ent.MemberRole, bool) {
	role := strings.TrimSpace(strings.ToLower(value))
	switch ent.MemberRole(role) {
	case ent.MemberRoleStudent, ent.MemberRoleTeacher, ent.MemberRoleAdmin, ent.MemberRoleStaff, ent.MemberRoleParent:
		return ent.MemberRole(role), true
	default:
		return "", false
	}
}

func containsRole(roles []ent.MemberRole, role ent.MemberRole) bool {
	for _, existing := range roles {
		if existing == role {
			return true
		}
	}

	return false
}

func orderRolesWithPrimary(roles []ent.MemberRole, role ent.MemberRole) []ent.MemberRole {
	ordered := make([]ent.MemberRole, 0, len(roles))
	ordered = append(ordered, role)

	for _, existing := range roles {
		if existing == role {
			continue
		}
		ordered = append(ordered, existing)
	}

	return ordered
}

func normalizeKnownRoles(roles []ent.MemberRole) []ent.MemberRole {
	normalized := make([]ent.MemberRole, 0, len(roles))

	for _, role := range roles {
		parsedRole, ok := parseKnownMemberRole(string(role))
		if !ok || containsRole(normalized, parsedRole) {
			continue
		}
		normalized = append(normalized, parsedRole)
	}

	return normalized
}
