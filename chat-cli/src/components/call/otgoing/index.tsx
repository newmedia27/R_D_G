import { useCallStore } from "@/stores/call"
import { useWsStore } from "@/stores/ws"
import { useShallow } from "zustand/shallow"

export const OutgoingCall = () => {
	const { rejectCall, roomId } = useCallStore(
		useShallow((state) => ({
			rejectCall: state.rejectCall,
			roomId: state.roomId,
		}))
	)
	function handleReject() {
		rejectCall()
		useWsStore.getState().send({
			type: "call.reject",
			room_id: roomId,
		})
	}
	return (
		<div>
			<p>Очікування відповіді...</p>
			<button onClick={handleReject}>Скасувати</button>
		</div>
	)
}
