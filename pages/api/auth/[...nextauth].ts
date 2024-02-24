import NextAuth from "next-auth"
import { OAuthConfig } from "next-auth/providers"

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
        role: profile.role,
      }),
      clientId: "test",
    } satisfies OAuthConfig<{
      sub: string
      name: string
      role: string
    }>,
  ],
  callbacks: {
    session: ({ session, token }) => ({
      ...session,
      user: {
        id: token.id,
        name: token.name,
        role: token.role,
      },
    }),
    jwt: ({ profile, token }) =>
      profile
        ? {
            ...token,
            id: profile.sub,
            username: profile.name,
            role: profile.role,
          }
        : token,
  },
})
