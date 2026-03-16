import { Input } from "@/components/ui/form/input/input"
import styles from "./index.module.sass"
import classNames from "classnames"
import { useEffect } from "react"
import { useChatStore } from "@/stores/chat"
import { useShallow } from "zustand/shallow"
import { AddChatModal } from "../addChatModal"
import { ChatItem } from "../chat-item"
import { useNavigate, useParams } from "react-router-dom"
import { useAuthStore } from "@/stores/auth"

export default function Sidebar() {
	const { id } = useParams()
	const { getChats, chats, members } = useChatStore(
		useShallow((state) => ({
			chats: state.chats,
			members: state.members,
			activeChat: state.activeChat,
			getChats: state.getChats,
		}))
	)
	const currentUserId = useAuthStore((state) => state.ownerId)
	const navigate = useNavigate()

	useEffect(() => {
		getChats()
	}, [getChats])

	return (
		<aside className={styles.wrapper}>
			<header className={styles.header}>
				{currentUserId
					? members?.[currentUserId]?.username
					: "username"}
			</header>
			<div className={classNames(styles.search, styles.content)}>
				<Input
					className={styles.search__input}
					placeholder="Search..."
				/>
			</div>
			<div
				className={classNames(styles.title__container, styles.content)}
			>
				<h3 className={styles.title}>Messages</h3>
				<AddChatModal>
					<button className={styles.add__chat} type="button">
						+
					</button>
				</AddChatModal>
			</div>

			<div className={classNames(styles.content)}>
				{chats.map((c) => {
					const memberId = c.members.find(
						(id) => id !== currentUserId
					)

					return (
						<ChatItem
							key={c.id}
							onClick={() => {
								navigate(`/chats/${c.id}`)
							}}
							chat={c}
							member={
								c.type === "private"
									? members?.[memberId as string]
									: undefined
							}
							active={c.id === id}
						/>
					)
				})}
			</div>
		</aside>
	)
}
