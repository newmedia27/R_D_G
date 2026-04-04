import { signOutRequest } from "@/api/auth/auth"
import { create } from "zustand"
import { devtools, persist } from "zustand/middleware"

type Token = string | null
interface Auth {
	token: Token
	isAuth: boolean
	ownerId: string | null

	setToken: (token: Token, usrID: string) => void
	setTokenAfterRefresh: (token: Token) => void
	resetToken: () => void
	resetTokenByRefresh: () => void
	signOut: () => Promise<void>
}

export const useAuthStore = create<Auth>()(
	persist(
		devtools((set) => ({
			token: null,
			isAuth: false,
			ownerId: null,

			setToken: (token: Token, usrID: string) =>
				set(() => ({
					token,
					isAuth: !!token,
					ownerId: usrID,
				})),
			setTokenAfterRefresh: (token: Token) =>
				set(() => ({
					token,
					isAuth: !!token,
				})),
			resetToken: () =>
				set(() => ({
					token: null,
					isAuth: false,
					ownerId: null,
				})),
			resetTokenByRefresh: () => ({
				token: null,
			}),
			signOut: async () => {
				try {
					await signOutRequest()
				} catch (err) {
					console.log("err", err)
				} finally {
					localStorage.removeItem("auth-storage")
				}
			},
		})),
		{ name: "auth-storage" }
	)
)

export const setTokenSelector = (state: Auth) => state.setToken
export const isAuthSelector = (state: Auth) => state.isAuth
