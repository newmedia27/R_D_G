export type User = {
	id: string
	username: string
	email: string
	created_at: string
	updated_at: string
}

export type UserMap = Record<string, User>
