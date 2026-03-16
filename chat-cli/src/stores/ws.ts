import { create } from "zustand"
import { devtools } from "zustand/middleware"
import { useMessageStore } from "./message"
import { useAuthStore } from "./auth"
import { useChatStore } from "./chat"
import type { Message } from "@/components/message.types"
import { refreshToken } from "@/api"
import type { Chat } from "@/components/chat.types"
import { useCallStore } from "./call"

const WS_URL = import.meta.env.VITE_WS_URL ?? "ws://localhost:8000/api/ws"

interface wsStore {
	ws: WebSocket | null
	connected: boolean

	connect: (token: string) => void
	disconnect: () => void
	send: (event: object) => void
}

export const useWsStore = create<wsStore>()(
	devtools((set, get) => ({
		ws: null,
		connected: false,

		connect: (token) => {
			const ws = new WebSocket(`${WS_URL}?token=${token}`)

			ws.onopen = () => {
				set({ connected: true })
				console.log("ws connected!!")
			}

			ws.onclose = (e) => {
				if (get().ws === ws) {
					console.log("ws closed", e.code, e.reason, e.wasClean)
					set({ connected: false, ws: null })
				}
			}

			ws.onerror = (e) => {
				console.error("ws error", e)
			}

			ws.onmessage = (e) => {
				try {
					const event = JSON.parse(e.data)
					handleEvent(event)
				} catch (err) {
					console.error("ws parse error", err)
				}
			}

			set({ ws })
		},

		disconnect: () => {
			get().ws?.close()
			set({ ws: null, connected: false })
		},

		send: (event) => {
			const { ws, connected } = get()
			if (!ws || !connected) return
			ws.send(JSON.stringify(event))
		},
	}))
)

function handleEvent(event: { type: string; [key: string]: unknown }) {
	switch (event.type) {
		case "message.new":
			useMessageStore.getState().addMessage(event.message as Message)
			useChatStore
				.getState()
				.updateLastMessage(
					event.chat_id as string,
					event.message as Message
				)
			break
		case "message.edited":
			useMessageStore.getState().editMessage(event.message as Message)
			break
		case "message.deleted":
			useMessageStore
				.getState()
				.removeMessage(
					event.message_id as string,
					event.chat_id as string
				)
			break
		case "auth.expiring_soon":
			console.log("exp")
			handleTokenExpiry()
			break
		case "chat.created":
			useChatStore.getState().addChat(event.chat as Chat)
			break
		case "call.invite":
			useCallStore
				.getState()
				.incomingCall(event.room_id as string, event.user_id as string)
			break
		case "call.accept": {
			const roomId = useCallStore.getState().roomId
			useCallStore.getState().acceptCall(roomId!)
			break
		}
		case "call.reject":
			useCallStore.getState().rejectCall()
			break
		case "call.end":
			useCallStore.getState().endCall()
			break
		default:
			console.warn("unknown ws event", event.type)
	}
}

async function handleTokenExpiry() {
	const oldToken = useAuthStore.getState().token
	if (!oldToken) return

	const token = await refreshToken(oldToken)

	if (!token) {
		useAuthStore.getState().signOut()
		return
	}

	useWsStore.getState().disconnect()
	useWsStore.getState().connect(token)
}
