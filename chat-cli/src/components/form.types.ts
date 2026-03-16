export interface ISignUp {
	email: string
	password: string
	passwordConfirm: string
	username: string
}
export type AuthResponse = {
	access_token: string
	user: {
		id: string
		email: string
		username: string
	}
}

export interface ISignIn {
	email: string
	password: string
}
