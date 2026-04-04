import { useEffect, useMemo, useState } from "react"
import styles from "./index.module.sass"
import { useDebounce } from "@/components/hooks/useDebounce"
import { useNavigate } from "react-router-dom"
import { useUserStore } from "@/stores/user"
import { useShallow } from "zustand/shallow"
import { useChatStore } from "@/stores/chat"
import { UserItem } from "../../user-item"
import type { ChatRequest } from "@/components/chat.types"

interface Props {
	onClose: () => void
}
export function Group({ onClose }: Props) {
	const [search, setSearch] = useState("")
	const [selectedUserIds, setSelectedUserIds] = useState<string[]>([])
	const [groupName, setGroupName] = useState("")
	const debounced = useDebounce(search)
	const navigate = useNavigate()
	const { searchUsers, searchedUsers, resetSerchedUsers } = useUserStore(
		useShallow((state) => ({
			searchUsers: state.searchUsers,
			resetSerchedUsers: state.resetSerchedUsers,
			searchedUsers: state.searchedUsers,
		}))
	)
	const { createGroupChat, members } = useChatStore(
		useShallow((state) => ({
			createGroupChat: state.createGroupChat,
			members: state.members,
		}))
	)

	useEffect(() => {
		if (debounced.trim().length < 3) {
			return resetSerchedUsers()
		}
		searchUsers(debounced)
	}, [debounced, resetSerchedUsers, searchUsers])

	const users = useMemo(() => {
		return searchedUsers?.filter((s) => !selectedUserIds?.includes(s.id))
	}, [searchedUsers, selectedUserIds])

	function handleClose() {
		resetSerchedUsers()
		onClose()
	}
	async function handleCreateChat() {
		if (selectedUserIds?.length) {
			const req: ChatRequest = {
				name: groupName,
				members: selectedUserIds,
			}
			const chat = await createGroupChat(req)
			if (chat) {
				navigate(`/chats/${chat.id}`)
			}
		}
		handleClose()
	}

	function toggleUser(id: string) {
		setSelectedUserIds((s) =>
			s.includes(id) ? s.filter((f) => f !== id) : [...s, id]
		)
	}
	return (
		<>
			<input
				className={styles.input}
				placeholder="Назва групи..."
				value={groupName}
				onChange={(e) => setGroupName(e.target.value)}
				autoFocus
			/>

			<input
				className={styles.input}
				placeholder="Пошук учасників..."
				value={search}
				onChange={(e) => setSearch(e.target.value)}
			/>

			{selectedUserIds.length > 0 && (
				<div className={styles.chips}>
					{selectedUserIds.map((id) => (
						<div key={id} className={styles.chip}>
							{members?.[id]?.username || id}
							<span
								className={styles.chipRemove}
								onClick={() => toggleUser(id)}
							>
								✕
							</span>
						</div>
					))}
				</div>
			)}

			<div className={styles.userList}>
				{users?.map((u) => {
					return (
						<UserItem
							key={u.id}
							email={u.email}
							username={u.username}
							avatar={u.username.charAt(0).toUpperCase()}
							onClick={() => {
								toggleUser(u.id)
							}}
						/>
					)
				})}
			</div>

			<div className={styles.footer}>
				<button
					className={styles.submitBtn}
					disabled={!groupName.trim() || selectedUserIds.length === 0}
					onClick={handleCreateChat}
				>
					Створити групу
					{selectedUserIds.length > 0 &&
						`(${selectedUserIds.length})`}
				</button>
			</div>
		</>
	)
}
