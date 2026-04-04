import { useChatStore } from "@/stores/chat"
import styles from "./index.module.sass"
import type { Message } from "@/components/message.types"
import { format, isSameDay } from "date-fns"
import { Fragment, useEffect, useRef } from "react"
import { uk } from "date-fns/locale"

const showDateSeparator = (current: Message, prev: Message | null) => {
	if (!prev) return true
	return !isSameDay(new Date(current.created_at), new Date(prev.created_at))
}

const DateSeparator = ({ date }: { date: string }) => {
	return (
		<div className={styles.dateSeparator}>
			<span className={styles.dateSeparatorLine} />
			<span className={styles.dateSeparatorText}>
				{format(new Date(date), "d MMMM yyyy 'р.'", { locale: uk })}
			</span>
			<span className={styles.dateSeparatorLine} />
		</div>
	)
}

interface Props {
	messages: Message[]
	currentUserId: string | null
}

export const MessageList = ({ messages, currentUserId }: Props) => {
	const bottomRef = useRef<HTMLDivElement>(null)

	useEffect(() => {
		bottomRef.current?.scrollIntoView({ behavior: "smooth" })
	}, [messages])
	if (!messages?.length) {
		return <div>Nothing messages yet...</div>
	}
	return (
		<div className={styles.list}>
			{messages.map((msg, i) => {
				const isOwn = msg.user_id === currentUserId
				const showSender =
					!isOwn &&
					(i === 0 || messages[i - 1].user_id !== msg.user_id)
				return (
					<Fragment key={msg.id}>
						{showDateSeparator(msg, messages[i - 1] ?? null) && (
							<DateSeparator date={msg.created_at} />
						)}
						<MessageItem
							message={msg}
							isOwn={isOwn}
							showSender={showSender}
						/>
					</Fragment>
				)
			})}
			<div ref={bottomRef} />
		</div>
	)
}

interface MessageItemProps {
	message: Message
	isOwn: boolean
	showSender: boolean
}

const MessageItem = ({ message, isOwn, showSender }: MessageItemProps) => {
	const members = useChatStore((state) => state.members)
	return (
		<div className={`${styles.row} ${isOwn ? styles.rowOwn : ""}`}>
			{!isOwn && (
				<div
					className={`${styles.avatar} ${showSender ? "" : styles.avatarHidden}`}
				>
					{members?.[message?.user_id]?.username
						.charAt(0)
						.toUpperCase()}
				</div>
			)}
			<div className={styles.content}>
				{showSender && !isOwn && (
					<span className={styles.sender}>
						{members?.[message?.user_id]?.username}
					</span>
				)}
				<div
					className={`${styles.bubble} ${isOwn ? styles.bubbleOwn : styles.bubbleOther}`}
				>
					<span className={styles.text}>{message.text}</span>
					<span className={styles.time}>
						{format(new Date(message.created_at), "HH:mm")}
					</span>
				</div>
			</div>
		</div>
	)
}
