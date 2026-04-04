import { useEffect, useState } from "react"
import styles from "./index.module.sass"
import { useDebounce } from "@/components/hooks/useDebounce"
import { useUserStore } from "@/stores/user"
import { UserItem } from "../../user-item"
import { useShallow } from "zustand/shallow"
import { useChatStore } from "@/stores/chat"
import { useNavigate } from "react-router-dom"

interface Props {
	onClose: () => void
}

export function Dirrect({ onClose }: Props) {
	const [search, setSearch] = useState("")
	const [selectedId, setSelectedId] = useState<string | null>(null)
	const debounced = useDebounce(search)
	const navigate = useNavigate()
	const { searchUsers, searchedUsers, resetSerchedUsers } = useUserStore(
		useShallow((state) => ({
			searchUsers: state.searchUsers,
			resetSerchedUsers: state.resetSerchedUsers,
			searchedUsers: state.searchedUsers,
		}))
	)
	const createChat = useChatStore((state) => state.createChat)

	useEffect(() => {
		if (debounced.trim().length < 3) {
			return resetSerchedUsers()
		}
		searchUsers(debounced)
	}, [debounced, resetSerchedUsers, searchUsers])

	function handleClose() {
		resetSerchedUsers()
		onClose()
	}

	async function handleCreateChat() {
		if (selectedId) {
			const chat = await createChat(selectedId)
			if (chat) {
				navigate(`/chats/${chat.id}`)
			}
		}
		handleClose()
	}
	return (
		<>
			<div className={styles.search}>
				<input
					className={styles.searchInput}
					placeholder="Search by username or email..."
					value={search}
					onChange={(e) => setSearch(e.target.value)}
					autoFocus
				/>
			</div>

			<div className={styles.userList}>
				{searchedUsers?.map((u) => (
					<UserItem
						key={u.id}
						email={u.email}
						username={u.username}
						selected={u.id === selectedId}
						avatar={u.username.charAt(0).toUpperCase()}
						onClick={() => {
							setSelectedId(u.id)
						}}
					/>
				))}
			</div>

			<div className={styles.footer}>
				<button
					className={styles.submitBtn}
					disabled={!selectedId}
					onClick={handleCreateChat}
				>
					Розпочати чат
				</button>
			</div>
		</>
	)
}
