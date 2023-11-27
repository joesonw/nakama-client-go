package nakama_client_go

import (
	"context"
	"fmt"
	"time"

	"github.com/heroiclabs/nakama-common/api"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/joesonw/nakama-client-go/httputil"
)

type Session struct {
	apiSession         *api.Session
	client             *Client
	expiresAt          time.Time
	autoRefreshSession bool
	vars               map[string]string
	userID             string
	username           string
	refreshExpiresAt   time.Time
}

func newSession(apiSession *api.Session, client *Client) (*Session, error) {
	s := &Session{
		apiSession:         apiSession,
		client:             client,
		autoRefreshSession: true,
		vars:               map[string]string{},
	}
	if err := s.decode(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Session) decode() error {
	var token struct {
		Vrs      map[string]string `json:"vrs"`
		UID      string            `json:"uid"`
		Username string            `json:"usn"`
		Exp      int64             `json:"exp"`
	}
	var refreshToken struct {
		Exp int64 `json:"exp"`
	}
	if err := unmarshalJWTBody(s.apiSession.Token, &token); err != nil {
		return err
	}
	if err := unmarshalJWTBody(s.apiSession.RefreshToken, &refreshToken); err != nil {
		return err
	}
	s.vars = token.Vrs
	s.userID = token.UID
	s.username = token.Username
	s.expiresAt = time.Unix(token.Exp, 0)
	s.refreshExpiresAt = time.Unix(refreshToken.Exp, 0)
	return nil
}

func (s *Session) Token() string {
	return s.apiSession.GetToken()
}

func (s *Session) Var(name string) string {
	return s.vars[name]
}

func (s *Session) IsExpired() bool {
	return s.expiresAt.Before(time.Now())
}

func (s *Session) IsRefreshExpired() bool {
	return s.refreshExpiresAt.Before(time.Now())
}

func (s *Session) UserID() string {
	return s.userID
}

func (s *Session) Username() string {
	return s.username
}

func (s *Session) DisableAutoRefresh() *Session {
	s.autoRefreshSession = false
	return s
}

func (s *Session) getUpToDatedToken(ctx context.Context) (*api.Session, error) {
	if s.IsRefreshExpired() {
		return nil, fmt.Errorf("refresh token expired")
	}
	if s.autoRefreshSession && s.apiSession.RefreshToken != "" && s.IsExpired() {
		if err := s.Refresh(ctx); err != nil {
			return nil, err
		}
	}
	return s.apiSession, nil
}

func (s *Session) attachSessionContext(ctx context.Context) (context.Context, error) {
	token, err := s.getUpToDatedToken(ctx)
	if err != nil {
		return nil, err
	}
	return httputil.WithBearerJWT(ctx, token.Token), nil
}

func (s *Session) Logout(ctx context.Context) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.SessionLogout(ctx, &api.SessionLogoutRequest{
		Token:        s.apiSession.Token,
		RefreshToken: s.apiSession.RefreshToken,
	})
	return err
}

func (s *Session) Refresh(ctx context.Context) error {
	res, err := s.client.api.SessionRefresh(httputil.WithBasicAuth(ctx, s.client.serverKey, ""), &api.SessionRefreshRequest{
		Token: s.apiSession.RefreshToken,
	})
	if err != nil {
		return err
	}
	s.apiSession = res
	return err
}

func (s *Session) AddGroupUsers(ctx context.Context, req *api.AddGroupUsersRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.AddGroupUsers(ctx, req)
	return err
}

func (s *Session) AddFriends(ctx context.Context, req *api.AddFriendsRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.AddFriends(ctx, req)
	return err
}

func (s *Session) BanGroupUsers(ctx context.Context, req *api.BanGroupUsersRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.BanGroupUsers(ctx, req)
	return err
}

func (s *Session) BlockFriends(ctx context.Context, req *api.BlockFriendsRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.BlockFriends(ctx, req)
	return err
}

func (s *Session) CreateGroup(ctx context.Context, req *api.CreateGroupRequest) (*api.Group, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.CreateGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) DeleteFriends(ctx context.Context, req *api.DeleteFriendsRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.DeleteFriends(ctx, req)
	return err
}

func (s *Session) DeleteGroup(ctx context.Context, req *api.DeleteGroupRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.DeleteGroup(ctx, req)
	return err
}

func (s *Session) DeleteLeaderboardRecord(ctx context.Context, req *api.DeleteLeaderboardRecordRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.DeleteLeaderboardRecord(ctx, req)
	return err
}

func (s *Session) DeleteNotifications(ctx context.Context, req *api.DeleteNotificationsRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.DeleteNotifications(ctx, req)
	return err
}

func (s *Session) DeleteStorageObjects(ctx context.Context, req *api.DeleteStorageObjectsRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.DeleteStorageObjects(ctx, req)
	return err
}

func (s *Session) DemoteGroupUsers(ctx context.Context, req *api.DemoteGroupUsersRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.DemoteGroupUsers(ctx, req)
	return err
}

func (s *Session) EmitEvent(ctx context.Context, req *api.Event) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.Event(ctx, req)
	return err
}

func (s *Session) GetAccount(ctx context.Context) (*api.Account, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.GetAccount(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) GetUsers(ctx context.Context, req *api.GetUsersRequest) ([]*api.User, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.GetUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.GetUsers(), nil
}

func (s *Session) JoinGroup(ctx context.Context, req *api.JoinGroupRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.JoinGroup(ctx, req)
	return err
}

func (s *Session) JoinTournament(ctx context.Context, req *api.JoinTournamentRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.JoinTournament(ctx, req)
	return err
}

func (s *Session) KickGroupUsers(ctx context.Context, req *api.KickGroupUsersRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.KickGroupUsers(ctx, req)
	return err
}

func (s *Session) LeaveGroup(ctx context.Context, req *api.LeaveGroupRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.LeaveGroup(ctx, req)
	return err
}

func (s *Session) ListChannelMessages(ctx context.Context, req *api.ListChannelMessagesRequest) (*api.ChannelMessageList, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ListChannelMessages(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) ListGroupUsers(ctx context.Context, req *api.ListGroupUsersRequest) (*api.GroupUserList, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ListGroupUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) ListUserGroups(ctx context.Context, req *api.ListUserGroupsRequest) (*api.UserGroupList, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ListUserGroups(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) ListGroups(ctx context.Context, req *api.ListGroupsRequest) (*api.GroupList, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ListGroups(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) LinkApple(ctx context.Context, req *api.AccountApple) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.LinkApple(ctx, req)
	return err
}

func (s *Session) LinkCustom(ctx context.Context, req *api.AccountCustom) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.LinkCustom(ctx, req)
	return err
}

func (s *Session) LinkDevice(ctx context.Context, req *api.AccountDevice) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.LinkDevice(ctx, req)
	return err
}

func (s *Session) LinkEmail(ctx context.Context, req *api.AccountEmail) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.LinkEmail(ctx, req)
	return err
}

func (s *Session) LinkFacebook(ctx context.Context, req *api.LinkFacebookRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.LinkFacebook(ctx, req)
	return err
}

func (s *Session) LinkFacebookInstantGame(ctx context.Context, req *api.AccountFacebookInstantGame) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.LinkFacebookInstantGame(ctx, req)
	return err
}

func (s *Session) LinkGameCenter(ctx context.Context, req *api.AccountGameCenter) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.LinkGameCenter(ctx, req)
	return err
}

func (s *Session) LinkGoogle(ctx context.Context, req *api.AccountGoogle) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.LinkGoogle(ctx, req)
	return err
}

func (s *Session) LinkSteam(ctx context.Context, req *api.LinkSteamRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.LinkSteam(ctx, req)
	return err
}

func (s *Session) ListFriends(ctx context.Context, req *api.ListFriendsRequest) (*api.FriendList, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ListFriends(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) ListLeaderboardRecords(ctx context.Context, req *api.ListLeaderboardRecordsRequest) (*api.LeaderboardRecordList, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ListLeaderboardRecords(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) ListLeaderboardRecordsAroundOwner(ctx context.Context, req *api.ListLeaderboardRecordsAroundOwnerRequest) (*api.LeaderboardRecordList, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ListLeaderboardRecordsAroundOwner(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) ListMatches(ctx context.Context, req *api.ListMatchesRequest) (*api.MatchList, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ListMatches(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) ListNotifications(ctx context.Context, req *api.ListNotificationsRequest) (*api.NotificationList, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ListNotifications(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) ListStorageObjects(ctx context.Context, req *api.ListStorageObjectsRequest) (*api.StorageObjectList, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ListStorageObjects(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) ListTournaments(ctx context.Context, req *api.ListTournamentsRequest) (*api.TournamentList, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ListTournaments(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) ListTournamentRecords(ctx context.Context, req *api.ListTournamentRecordsRequest) (*api.TournamentRecordList, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ListTournamentRecords(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) ListTournamentRecordsAroundOwner(ctx context.Context, req *api.ListTournamentRecordsAroundOwnerRequest) (*api.TournamentRecordList, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ListTournamentRecordsAroundOwner(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) PromoteGroupUsers(ctx context.Context, req *api.PromoteGroupUsersRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.PromoteGroupUsers(ctx, req)
	return err
}

func (s *Session) ReadStorageObjects(ctx context.Context, req *api.ReadStorageObjectsRequest) (*api.StorageObjects, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ReadStorageObjects(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) RPC(ctx context.Context, id string, payload []byte) (*api.Rpc, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.RpcFunc(ctx, &api.Rpc{Id: id, Payload: string(payload)})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) UnlinkApple(ctx context.Context, req *api.AccountApple) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.UnlinkApple(ctx, req)
	return err
}

func (s *Session) UnlinkCustom(ctx context.Context, req *api.AccountCustom) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.UnlinkCustom(ctx, req)
	return err
}

func (s *Session) UnlinkDevice(ctx context.Context, req *api.AccountDevice) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.UnlinkDevice(ctx, req)
	return err
}

func (s *Session) UnlinkEmail(ctx context.Context, req *api.AccountEmail) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.UnlinkEmail(ctx, req)
	return err
}

func (s *Session) UnlinkFacebook(ctx context.Context, req *api.AccountFacebook) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.UnlinkFacebook(ctx, req)
	return err
}

func (s *Session) UnlinkFacebookInstantGame(ctx context.Context, req *api.AccountFacebookInstantGame) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.UnlinkFacebookInstantGame(ctx, req)
	return err
}

func (s *Session) UnlinkGameCenter(ctx context.Context, req *api.AccountGameCenter) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.UnlinkGameCenter(ctx, req)
	return err
}

func (s *Session) UnlinkGoogle(ctx context.Context, req *api.AccountGoogle) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.UnlinkGoogle(ctx, req)
	return err
}

func (s *Session) UnlinkSteam(ctx context.Context, req *api.AccountSteam) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.UnlinkSteam(ctx, req)
	return err
}

func (s *Session) UpdateAccount(ctx context.Context, req *api.UpdateAccountRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.UpdateAccount(ctx, req)
	return err
}

func (s *Session) UpdateGroup(ctx context.Context, req *api.UpdateGroupRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.UpdateGroup(ctx, req)
	return err
}

func (s *Session) ValidatePurchaseApple(ctx context.Context, req *api.ValidatePurchaseAppleRequest) (*api.ValidatePurchaseResponse, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ValidatePurchaseApple(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) ValidatePurchaseGoogle(ctx context.Context, req *api.ValidatePurchaseGoogleRequest) (*api.ValidatePurchaseResponse, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ValidatePurchaseGoogle(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) ValidatePurchaseHuawei(ctx context.Context, req *api.ValidatePurchaseHuaweiRequest) (*api.ValidatePurchaseResponse, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.ValidatePurchaseHuawei(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) WriteLeaderboardRecord(ctx context.Context, req *api.WriteLeaderboardRecordRequest) (*api.LeaderboardRecord, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.WriteLeaderboardRecord(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) WriteStorageObjects(ctx context.Context, req *api.WriteStorageObjectsRequest) (*api.StorageObjectAcks, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.WriteStorageObjects(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) WriteTournamentRecord(ctx context.Context, req *api.WriteTournamentRecordRequest) (*api.LeaderboardRecord, error) {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.client.api.WriteTournamentRecord(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Session) ImportSteamFriends(ctx context.Context, req *api.ImportSteamFriendsRequest) error {
	ctx, err := s.attachSessionContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.client.api.ImportSteamFriends(ctx, req)
	return err
}
