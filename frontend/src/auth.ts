import NextAuth from "next-auth"
import Google from "next-auth/providers/google"

export const { auth, handlers, signIn, signOut } = NextAuth({
  providers: [
    Google({
      clientId: process.env.AUTH_GOOGLE_ID,
      clientSecret: process.env.AUTH_GOOGLE_SECRET,
    })
  ],
  pages: {
    signIn: '/auth/signin',
    error: '/auth/signin', // Redirect errors to our signin page too
  },
  callbacks: {
    async signIn({ user, account, profile }) {
      if (account?.provider === "google") {
        try {
          // Call your backend to create or update user
          const response = await fetch(`${process.env.NEXTAUTH_URL || 'http://localhost:3000'}/api/auth/google`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              googleId: account.providerAccountId,
              email: user.email,
              name: user.name,
              picture: user.image,
              accessToken: account.access_token,
              refreshToken: account.refresh_token,
              expiresAt: account.expires_at ? new Date(account.expires_at * 1000).toISOString() : undefined,
            }),
          })

          if (response.ok) {
            const backendUser = await response.json()
            // Store backend user info in user object for JWT
            user.id = backendUser.id
            ;(user as any).sessionToken = backendUser.sessionToken
            return true
          } else {
            console.error('Failed to create/login user in backend')
            return false
          }
        } catch (error) {
          console.error('Error during sign in:', error)
          return false
        }
      }
      return true
    },
    async jwt({ token, user }) {
      // First time JWT callback is run, user object is available
      if (user) {
        token.id = user.id
        token.sessionToken = (user as any).sessionToken
      }
      return token
    },
    async session({ session, token }) {
      // Send properties to the client
      if (session.user) {
        session.user.id = token.id as string
        session.user.sessionToken = token.sessionToken as string
      }
      return session
    },
    async redirect({ url, baseUrl }) {
      // Always redirect to /home after successful signin
      if (url.startsWith("/") && !url.startsWith("/auth/")) return `${baseUrl}/home`
      // Allow callback URLs on the same origin
      if (new URL(url).origin === baseUrl) return url
      return `${baseUrl}/home`
    },
  },
  session: {
    strategy: "jwt",
    maxAge: 24 * 60 * 60, // 24 hours
  },
}) 