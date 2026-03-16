export const CALL_STATUSES = {
	IDLE: "idle",
	INCOMING: "incoming",
	OUTGOING: "outgoing",
	CONNECTED: "connected",
}
export type STATUS =
	| typeof CALL_STATUSES.IDLE
	| typeof CALL_STATUSES.INCOMING
	| typeof CALL_STATUSES.OUTGOING
	| typeof CALL_STATUSES.CONNECTED

export interface CallTokenResponse {
	token: string
	url: string
}

export interface StartCallRequest {
	chat_id: string
}
