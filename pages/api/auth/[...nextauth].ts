import NextAuth from "next-auth";

export default NextAuth({
  providers: [
    {
      id: "go-oidc-idp-example",
      name: "Go OIDC IdP Example",
      type: "oauth",
      wellKnown: "http://localhost:8080/.well-known/openid-configuration",
      authorization: { params: { response_type: "id_token" } },
      idToken: true,
      checks: ["nonce"],
      profile: (profile) => ({
        id: profile.sub,
        name: profile.name,
      }),
      clientId: "test",
    },
  ],
});
