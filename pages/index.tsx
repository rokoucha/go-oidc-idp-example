import { signIn, signOut, useSession } from "next-auth/react";
import React from "react";

export default function Component() {
  const { data: session } = useSession();

  return session ? (
    <>
      Signed in as {session.user?.name} <br />
      <button onClick={() => signOut()}>Sign out</button>
    </>
  ) : (
    <>
      Not signed in <br />
      <button onClick={() => signIn()}>Sign in</button>
    </>
  );
}
