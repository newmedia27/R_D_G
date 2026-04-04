import styles from "./index.module.sass"
import { createPortal } from "react-dom"
import Draggable from "react-draggable"
import { CALL_STATUSES } from "../call.types"
import { useCallStore } from "@/stores/call"
import { IncomingCall } from "./incoming"
import { OutgoingCall } from "./otgoing"
import { ActiveCall } from "./active"
import classNames from "classnames"

export function CallModal({ open }: { open: boolean }) {
	const { status } = useCallStore()

	if (status === CALL_STATUSES.IDLE || !open) return null

	const isConnected = status === CALL_STATUSES.CONNECTED

	return createPortal(
		<div
			className={classNames(styles.modal, {
				[styles.modalConnected]: isConnected,
			})}
		>
			{status === "incoming" && <IncomingCall />}
			{status === "outgoing" && <OutgoingCall />}
			{status === "connected" && <ActiveCall />}
		</div>,
		document.body
	)
}
