import { useCallStore } from "@/stores/call"
import styles from "./index.module.sass"
import { useChatStore } from "@/stores/chat"
import { useShallow } from "zustand/shallow"
import { useWsStore } from "@/stores/ws"

export const IncomingCall = () => {
	const { roomId, userIds, acceptCall, rejectCall } = useCallStore(
		useShallow((state) => ({
			roomId: state.roomId,
			userIds: state.userIds,
			acceptCall: state.acceptCall,
			rejectCall: state.rejectCall,
		}))
	)
	const members = useChatStore((state) => state.members)
	const caller = userIds?.[0] ? members?.[userIds[0]] : null

	function handleAccept() {
		acceptCall(roomId!)
		useWsStore.getState().send({
			type: "call.accept",
			room_id: roomId,
		})
	}

	function handleReject() {
		rejectCall()
		useWsStore.getState().send({
			type: "call.reject",
			room_id: roomId,
		})
	}

	return (
		<div className={styles.incoming}>
			<div className={styles.avatar}>
				{caller?.username?.charAt(0).toUpperCase()}
			</div>
			<p className={styles.callerName}>{caller?.username}</p>
			<p className={styles.callerStatus}>Вхідний дзвінок...</p>
			<div className={styles.actions}>
				<button className={styles.rejectBtn} onClick={handleReject}>
					✕
				</button>
				<button className={styles.acceptBtn} onClick={handleAccept}>
					✓
				</button>
			</div>
		</div>
	)
}
