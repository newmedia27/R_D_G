export type Chat = {
	created_at: string
	description: string
	id: string
	last_message: LastMessage | null
	members: Members
	name: string
	owner_id: string
	type: ChatType
	updated_at: string
}

type LastMessage = {
	created_at: string
	text: string
	user_id: string
}

type Members = string[]

type ChatType = "private" | "group"

export interface ChatRequest {
	description?: string
	members: Members
	name: string
}
