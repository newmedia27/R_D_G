import { getProfileRequest, searchUsersRequest } from "@/api/user/user"
import type { User, UserMap } from "@/components/user.types"
import { create } from "zustand"
import { devtools } from "zustand/middleware"
import { useChatStore } from "./chat"

interface UserStore {
	profile: User
	searchedUsers: User[] | null

	getProfile: () => Promise<void>
	searchUsers: (search: string) => Promise<void>
	resetSerchedUsers: () => void
}

export const useUserStore = create<UserStore>()(
	devtools((set) => ({
		profile: null,
		searchedUsers: null,

		searchUsers: async (search: string) => {
			try {
				const { data } = await searchUsersRequest(search)
				set({ searchedUsers: data.users })
				const { members, setMembers } = useChatStore.getState()
				const users = data.users.reduce((acc: UserMap, u) => {
					if (!members?.[u.id]) {
						acc[u.id] = u
					}
					return acc
				}, {})
				setMembers({ ...members, ...users })
			} catch (err) {
				console.log("err", err)
			}
		},
		resetSerchedUsers: () => {
			set({ searchedUsers: null })
		},
		getProfile: async () => {
			try {
				const { data } = await getProfileRequest()
				set({ profile: data?.user })
			} catch (err) {
				console.log("err", err)
			}
		},
	}))
)
