export type Message = {
	id: string
	chat_id: string
	user_id: string
	type: MessageType
	text: string
	file_id: string
	isEdited: boolean
	isDeleted: boolean
	created_at: string
	updated_at: string
}

type MessageType = "text" | "file"
