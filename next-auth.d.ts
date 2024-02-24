import "next-auth"

declare module "next-auth" {
  interface Profile {
    id: string
    name: string
    role: string
  }

  interface User {
    id: string
    name: string
    role: string
  }

  interface Session {
    user: {
      id: string
      name: string
      role: string
    }
  }
}

declare module "next-auth/jwt" {
  interface JWT {
    id: string
    name: string
    role: string
  }
}
