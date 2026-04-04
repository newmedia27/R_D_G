import { getMessagesRequest } from "@/api/messages/messages"
import type { Message } from "@/components/message.types"
import { create } from "zustand"
import { devtools } from "zustand/middleware"
import { useWsStore } from "./ws"

interface MessageStore {
	messages: Record<string, Message[]>

	getMessages: (chatId: string) => Promise<void>
	addMessage: (message: Message) => Promise<void>
	editMessage: (message: Message) => Promise<void>
	removeMessage: (id: string, chatId: string) => Promise<void>
	sendMessage: (chatId: string, text: string) => void
}

export const useMessageStore = create<MessageStore>()(
	devtools((set) => ({
		messages: {},

		getMessages: async (chatId) => {
			try {
				const { data } = await getMessagesRequest(chatId)
				set((state) => ({
					messages: {
						...state?.messages,
						[chatId]: data,
					},
				}))
			} catch (err) {
				console.log("err", err)
			}
		},
		sendMessage: (chatId, text) => {
			const ws = useWsStore.getState().ws
			if (!ws) return
			ws.send(
				JSON.stringify({
					type: "message.send",
					chat_id: chatId,
					text,
				})
			)
		},
		addMessage: (message: Message) => {
			set((state) => ({
				messages: {
					...state.messages,
					[message.chat_id]: [
						...(state.messages[message.chat_id] ?? []),
						message,
					],
				},
			}))
		},

		editMessage: (message: Message) => {
			set((state) => ({
				messages: {
					...state.messages,
					[message.chat_id]: (
						state.messages[message.chat_id] ?? []
					).map((m) => (m.id === message.id ? message : m)),
				},
			}))
		},

		removeMessage: (id: string, chatId: string) => {
			set((state) => ({
				messages: {
					...state.messages,
					[chatId]: (state.messages[chatId] ?? []).filter(
						(m) => m.id !== id
					),
				},
			}))
		},
	}))
)
