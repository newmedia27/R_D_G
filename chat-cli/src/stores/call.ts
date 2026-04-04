import { startCallRequest } from "@/api/call/call"
import { CALL_STATUSES, type STATUS } from "@/components/call.types"
import { create } from "zustand"
import { devtools } from "zustand/middleware"

interface CallStore {
	callModal: boolean
	roomId: string | null
	userIds: string[] | null
	status: STATUS
	token: string | null
	url: string | null

	startCall: (chatId: string, userIds: string[]) => Promise<void>
	incomingCall: (roomId: string, userId: string) => void
	acceptCall: (roomId: string) => Promise<void>
	rejectCall: () => void
	endCall: () => void
}

const initialState = {
	callModal: false,
	roomId: null,
	userIds: null,
	status: CALL_STATUSES.IDLE as STATUS,
	token: null,
	url: null,
}

export const useCallStore = create<CallStore>()(
	devtools((set) => ({
		callModal: false,
		roomId: null,
		userIds: null,
		status: "idle",
		token: null,
		url: null,

		startCall: async (chatId, userIds) => {
			set({ roomId: chatId, userIds, status: CALL_STATUSES.OUTGOING })
			try {
				const { data } = await startCallRequest({
					chat_id: chatId as string,
				})
				console.log("data", data)
				set({ url: data?.url, token: data?.token, callModal: true })
			} catch (err) {
				console.log("err", err)
				set(initialState)
			}
		},
		incomingCall: (roomId, userId) => {
			set({
				status: CALL_STATUSES.INCOMING,
				roomId,
				userIds: [userId],
				callModal: true,
			})
		},
		acceptCall: async (chatId) => {
			try {
				const { data } = await startCallRequest({
					chat_id: chatId as string,
				})
				set({
					status: CALL_STATUSES.CONNECTED,
					url: data?.url,
					token: data?.token,
				})
			} catch (err) {
				console.log("err", err)
				set(initialState)
			}
		},
		rejectCall: () => {
			set(initialState)
		},
		endCall: () => {
			set(initialState)
		},
	}))
)
