import type { AxiosResponse } from "axios"
import { api } from ".."
import type { Message } from "@/components/message.types"

export async function getMessagesRequest(
	chatId: string
): Promise<AxiosResponse<Message[]>> {
	const url = `/chats/${chatId}/messages`
	return await api.get<Message[]>(url)
}
