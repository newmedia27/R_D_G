import { useParams } from "react-router-dom"
import styles from "./index.module.sass"
import { useEffect, useMemo } from "react"
import { useChatStore } from "@/stores/chat"
import { useShallow } from "zustand/shallow"
import { useAuthStore } from "@/stores/auth"
import type { User } from "@/components/user.types"
import { MessageInput } from "./input"
import { MessageList } from "./message-list"
import { useMessageStore } from "@/stores/message"
import { useCallStore } from "@/stores/call"
import { useWsStore } from "@/stores/ws"
import { CallModal } from "@/components/call"

export default function Main() {
	const { id } = useParams()
	const currentUserId = useAuthStore((state) => state.ownerId)
	const callModal = useCallStore((state) => state.callModal)

	const { members, activeChat, getChat } = useChatStore(
		useShallow((state) => ({
			chats: state.chats,
			members: state.members,
			getChat: state.getChat,
			activeChat: state.activeChat,
		}))
	)
	const { messages, getMessages, sendMessage } = useMessageStore(
		useShallow((state) => ({
			messages: state.messages,
			getMessages: state.getMessages,
			sendMessage: state.sendMessage,
		}))
	)

	const startCall = useCallStore((state) => state.startCall)
	const info = useMemo(() => {
		if (!activeChat) return null
		if (activeChat.type === "private") {
			const memberId = activeChat.members.find((m) => m !== currentUserId)
			return memberId ? members?.[memberId] : null
		}
		return { username: activeChat.name } as User
	}, [activeChat, members, currentUserId])

	useEffect(() => {
		if (id) {
			getChat({ chatId: id })
			getMessages(id)
		}
	}, [id, getChat, getMessages])

	if (!id) {
		return (
			<div className={styles.content}>Please pick or create chat!!</div>
		)
	}

	async function handleCall() {
		await startCall(id as string, Object.keys(members ?? {}))
		console.log("object")
		useWsStore.getState().send({
			type: "call.invite",
			room_id: id,
		})
	}

	return (
		<main className={styles.wrapper}>
			<header className={styles.header}>
				{info && <User info={info} />}
				<button type="button" onClick={handleCall}>
					<svg
						width="20"
						height="20"
						viewBox="0 0 20 20"
						fill="none"
						xmlns="http://www.w3.org/2000/svg"
					>
						<path
							d="M4.5 2C4.5 2 3 2 2 3.5C1 5 2 7.5 4 9.5C6 11.5 8.5 13 10 13.5C11.5 14 13.5 14.5 15 13.5C16.5 12.5 17.5 11 17.5 11L14.5 9L13 10.5C13 10.5 11.5 9.5 10 8C8.5 6.5 7.5 5 7.5 5L9 3.5L6.5 2H4.5Z"
							stroke="currentColor"
							strokeWidth="1.5"
							strokeLinejoin="round"
						/>
					</svg>
				</button>
			</header>
			<div className={styles.content}>
				<MessageList
					messages={messages[id]}
					currentUserId={currentUserId}
				/>
			</div>
			<div className={styles.message__input_wrapper}>
				<MessageInput
					placeholder="Enter message..."
					onSend={(text) => sendMessage(id, text)}
				/>
			</div>
			<CallModal open={callModal} />
		</main>
	)
}

interface UserProps {
	info: User
}

const User = ({ info }: UserProps) => {
	return (
		<div className={styles.user}>
			<div className={styles.avatarWrapper}>
				<div className={styles.avatar}>
					{info.username[0].toUpperCase()}
				</div>
			</div>
			<div className={styles.info}>
				<span className={styles.name}>{info.username}</span>
				<span className={styles.status}>
					<span className={styles.onlineDot} />
					онлайн
				</span>
			</div>
		</div>
	)
}
