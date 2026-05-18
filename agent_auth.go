package chronary

import (
	"context"
	"net/http"
)

// AgentAuthService handles agent self-signup operations.
//
// SignUp is unauthenticated and can be invoked on a client constructed with
// WithAnonymous() (no API key). Verify must be called on a client constructed
// with the restricted API key returned from a successful new-org SignUp.
type AgentAuthService struct {
	service
}

// SignUp registers a new agent + org. Sends an OTP to the provided email.
//
// Two response shapes:
//   - **New org** — IsNewOrg() returns true. APIKey is populated and is
//     restricted to /v1/agent/verify until the OTP is submitted.
//   - **Existing org** — IsNewOrg() returns false. Only Message is populated;
//     no credentials are leaked when the email matches an existing org
//     (enumeration defense).
//
// Returns a 409 ChronaryError with type "tos_version_stale" if TosVersion is
// not the currently-published version.
func (s *AgentAuthService) SignUp(ctx context.Context, params *AgentSignUpParams, opts ...RequestOption) (*AgentSignUpResponse, error) {
	var resp AgentSignUpResponse
	err := s.client.do(ctx, http.MethodPost, "/v1/agent/sign-up", params, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Verify submits the OTP from SignUp to unlock the API key.
//
// Must be called on a client constructed with the restricted APIKey returned
// from SignUp — NOT on the anonymous client used to issue SignUp.
//
// Returns a 400 ChronaryError when the OTP is wrong or expired.
func (s *AgentAuthService) Verify(ctx context.Context, params *AgentVerifyParams, opts ...RequestOption) (*AgentVerifyResponse, error) {
	var resp AgentVerifyResponse
	err := s.client.do(ctx, http.MethodPost, "/v1/agent/verify", params, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
