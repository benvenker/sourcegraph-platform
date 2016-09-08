import expect from "expect.js";
import * as React from "react";
import {AuthInfo, User} from "sourcegraph/api";
import {getChildContext, withUserContext} from "sourcegraph/app/user";
import * as UserActions from "sourcegraph/user/UserActions";
import {UserStore} from "sourcegraph/user/UserStore";
import {render} from "sourcegraph/util/testutil/renderTestUtils";

const sampleAuthInfo: AuthInfo = {UID: 1, Login: "u"} as AuthInfo;
const sampleUser: User = {UID: 1, Login: "u", Betas: [], BetaRegistered: false} as any as User;

const C = withUserContext((props) => null);
const renderAndGetContext = (c) => {
	const res = render(c, {});

	// Hack to get state, so we can pass it to getChildContext.
	const state = {};
	const e = new C({});
	e.reconcileState(state, {});
	e.onStateTransition(state, state);
	return Object.assign({}, res, {context: getChildContext(state)});
};

describe("withUserContext", () => {
	it("no accessToken", () => {
		UserStore.activeAccessToken = null;
		const res = renderAndGetContext(<C />);
		expect(res.actions).to.eql([]);
		expect(res.context).to.eql({authInfo: null, user: null, signedIn: false});
	});
	it("with accessToken, no authInfo yet", () => {
		UserStore.activeAccessToken = "t";
		const res = renderAndGetContext(<C />);
		expect(res.actions).to.eql([new UserActions.WantAuthInfo("t")]);
		expect(res.context).to.eql({authInfo: null, user: null, signedIn: true});
	});
	it("with accessToken, authInfo, and user", () => {
		UserStore.activeAccessToken = "t";
		UserStore.directDispatch(new UserActions.FetchedAuthInfo("t", sampleAuthInfo));
		UserStore.directDispatch(new UserActions.FetchedUser(1, sampleUser));
		const res = renderAndGetContext(<C />);
		expect(res.actions).to.eql([]);
		expect(res.context).to.eql({authInfo: {Login: "u", UID: 1}, user: sampleUser, signedIn: true});
	});
	it("with accessToken but empty authInfo object (indicating no user, expired accessToken, etc.)", () => {
		UserStore.activeAccessToken = "t";
		UserStore.directDispatch(new UserActions.FetchedAuthInfo("t", {} as any));
		const res = renderAndGetContext(<C />);
		expect(res.actions).to.eql([]);
		expect(res.context.signedIn).to.be(false);
	});
});
