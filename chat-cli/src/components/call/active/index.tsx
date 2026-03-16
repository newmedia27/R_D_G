import { useCallStore } from "@/stores/call"
import {
	LiveKitRoom,
	useLocalParticipant,
	ParticipantTile,
	useRoomContext,
	useTracks,
} from "@livekit/components-react"
import { Track } from "livekit-client"
import { useShallow } from "zustand/shallow"
import styles from "./index.module.sass"

const CallControls = ({ onLeave }: { onLeave: () => void }) => {
	const room = useRoomContext()
	const { localParticipant } = useLocalParticipant()

	const toggleMic = () =>
		localParticipant.setMicrophoneEnabled(
			!localParticipant.isMicrophoneEnabled
		)
	const toggleCamera = () =>
		localParticipant.setCameraEnabled(!localParticipant.isCameraEnabled)

	return (
		<div className={styles.controls}>
			<button onClick={toggleMic}>
				{localParticipant.isMicrophoneEnabled ? "🎤" : "🔇"}
			</button>
			<button onClick={toggleCamera}>
				{localParticipant.isCameraEnabled ? "📷" : "🚫"}
			</button>
			<button
				className={styles.leaveBtn}
				onClick={() => {
					room.disconnect()
					onLeave()
				}}
			>
				Leave
			</button>
		</div>
	)
}

const CallRoom = ({ onLeave }: { onLeave: () => void }) => {
	const tracks = useTracks([
		{ source: Track.Source.Camera, withPlaceholder: true },
		{ source: Track.Source.ScreenShare, withPlaceholder: false },
	])

	return (
		<div className={styles.room}>
			<div className={styles.grid}>
				{tracks.map((track) => (
					<ParticipantTile
						key={track.participant.identity}
						trackRef={track}
					/>
				))}
			</div>
			<CallControls onLeave={onLeave} />
		</div>
	)
}

export const ActiveCall = () => {
	const { token, url, endCall } = useCallStore(
		useShallow((state) => ({
			token: state.token,
			url: state.url,
			endCall: state.endCall,
		}))
	)

	if (!token || !url) return null

	return (
		<LiveKitRoom key={token} token={token} serverUrl={url} connect>
			<CallRoom onLeave={endCall} />
		</LiveKitRoom>
	)
}
