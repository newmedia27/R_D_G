import * as yup from "yup"

export const signUpSchema = yup.object().shape({
	email: yup.string().required().email(),
	username: yup.string().trim().required("required field").min(2).max(255),
	password: yup
		.string()
		.required()
		.matches(
			/^(?=.*\d)(?=.*[a-z])(?=.*[a-zA-Z]).{8,}$/,
			"Password must contain at least 8 characters, one letter and one number"
		),
	passwordConfirm: yup
		.string()
		.required()
		.oneOf([yup.ref("password")], "Passwords must match"),
})

export const signinSchema = yup.object().shape({
	email: yup.string().required().email(),
	password: yup
		.string()
		.required()
		.matches(
			/^(?=.*\d)(?=.*[a-z])(?=.*[a-zA-Z]).{8,}$/,
			"Password must contain at least 8 characters, one letter and one number"
		),
})
