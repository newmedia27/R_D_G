import { RouterProvider } from "react-router-dom"
import { BasicRouter } from "./routes"

function App() {
	return (
		<>
			<RouterProvider
				future={{
					v7_startTransition: true,
				}}
				router={BasicRouter}
			/>
		</>
	)
}

export default App
