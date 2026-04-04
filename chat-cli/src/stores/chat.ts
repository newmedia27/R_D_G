import {
	createDirrectChatRequest,
	createGroupChatRequest,
	getChatRequrest,
	getChatsRequest,
} from "@/api/chats/chats"
import { getUsersByIdsRequest } from "@/api/user/user"
import type { Chat, ChatRequest } from "@/components/chat.types"
import type { Message } from "@/components/message.types"
import type { UserMap } from "@/components/user.types"
import { create } from "zustand"
import { devtools } from "zustand/middleware"

type ChatStore = {
	chats: Chat[]
	members: UserMap | null
	activeChat: Chat | null

	getChats: () => Promise<void>
	getChat: ({ chatId }: { chatId: string }) => Promise<void>
	getMembers: (ids: string[]) => Promise<void>
	setMembers: (members: UserMap) => void
	createChat: (id: string) => Promise<Chat | null>
	createGroupChat: (values: ChatRequest) => Promise<Chat | null>
	updateLastMessage: (chatId: string, message: Message) => void
	addChat: (chat: Chat) => void
}

export const useChatStore = create<ChatStore>()(
	devtools((set, get) => ({
		chats: [],
		members: {},
		activeChat: null,

		getChat: async ({ chatId }) => {
			try {
				const { data } = await getChatRequrest({ chatId })
				set({ activeChat: data })
			} catch (err) {
				console.log("err", err)
			}
		},
		getChats: async () => {
			try {
				const { data } = await getChatsRequest()
				set({ chats: data })
				const arr = data.flatMap((c) => c.members)
				const ids = [...new Set(arr)]
				console.log("ids", ids)
				await get().getMembers(ids)
			} catch (err) {
				console.log("err", err)
			}
		},
		getMembers: async (ids) => {
			try {
				const { data } = await getUsersByIdsRequest(ids)
				set({ members: data.users })
			} catch (err) {
				console.log("err", err)
			}
		},
		setMembers: (members) => {
			set({ members })
		},
		createChat: async (id) => {
			try {
				const { data } = await createDirrectChatRequest(id)
				set((state) => ({
					activeChat: data,
					chats: [data, ...state.chats],
				}))
				return data
			} catch (err) {
				console.log("err", err)
				return null
			}
		},
		createGroupChat: async (values) => {
			try {
				const { data } = await createGroupChatRequest(values)
				set((state) => ({
					activeChat: data,
					chats: [data, ...state.chats],
				}))
				return data
			} catch (err) {
				console.log("err", err)
				return null
			}
		},
		updateLastMessage: (chatId: string, message: Message) => {
			set((state) => ({
				chats: state.chats.map((c) =>
					c.id === chatId ? { ...c, last_message: message } : c
				),
			}))
		},
		addChat: (chat) => {
			set((state) => ({ chats: [chat, ...state.chats] }))
		},
	}))
)
