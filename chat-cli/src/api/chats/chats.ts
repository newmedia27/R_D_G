import type { AxiosResponse } from "axios"
import { api } from ".."
import type { Chat, ChatRequest } from "@/components/chat.types"

export async function getChatsRequest(): Promise<AxiosResponse<Chat[]>> {
	const url = `/chats`
	return await api.get<Chat[]>(url)
}

export async function getChatRequrest({
	chatId,
}: {
	chatId: string
}): Promise<AxiosResponse<Chat>> {
	const url = `/chats/${chatId}`
	return await api.get<Chat>(url)
}

export async function createDirrectChatRequest(
	id: string
): Promise<AxiosResponse<Chat>> {
	const url = `/chats/private/${id}`
	return await api.post<Chat>(url)
}

export async function createGroupChatRequest(
	data: ChatRequest
): Promise<AxiosResponse<Chat>> {
	const url = `/chats/group`
	return await api.post<Chat>(url, data)
}
