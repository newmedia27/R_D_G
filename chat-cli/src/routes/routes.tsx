import {
	createBrowserRouter,
	createRoutesFromElements,
	Navigate,
	Outlet,
	Route,
} from "react-router-dom"
import AuthLayout from "@/components/auth/layout"
import SignIn from "@/components/auth/SignIn"
import SignUp from "@/components/auth/SignUp"
import { isAuthSelector, useAuthStore } from "@/stores/auth"
import { type ReactNode } from "react"
import ChatLayout from "@/components/chat/layout"
import ChatPage from "@/components/chat"

function ProtectedRoute() {
	const isAuth = useAuthStore(isAuthSelector)

	if (!isAuth) {
		return <Navigate to="/auth/signin" replace />
	}
	return <Outlet />
}

function WithAuthProtected({ children }: { children: Readonly<ReactNode> }) {
	const isAuth = useAuthStore(isAuthSelector)
	if (isAuth) {
		return <Navigate to="/" replace />
	}
	return <>{children}</>
}

export const BasicRouter = createBrowserRouter(
	createRoutesFromElements(
		<>
			<Route
				path="/auth/"
				element={
					<WithAuthProtected>
						<AuthLayout />
					</WithAuthProtected>
				}
			>
				<Route path="signin" element={<SignIn />} />
				<Route path="signup" element={<SignUp />} />
			</Route>
			<Route path="/" element={<ProtectedRoute />}>
				<Route index element={<Navigate to="/chats" />} />
				<Route path="chats" element={<ChatLayout />}>
					<Route index element={<ChatPage />} />
					<Route path=":id" element={<ChatPage />} />
				</Route>
			</Route>
			<Route path="*" element={<div>404 no found</div>} />
		</>
	)
)
