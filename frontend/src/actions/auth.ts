"use server"

import { signIn, signOut } from "@/auth"
import { AuthError } from "next-auth"

export async function signInWithGoogle() {
  try {
    await signIn("google", { redirectTo: "/home" })
  } catch (error) {
    if (error instanceof AuthError) {
      // Handle auth errors if needed
      throw error
    }
    throw error
  }
}

export async function signOutAction() {
  await signOut()
} 