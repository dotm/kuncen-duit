import React, { ReactElement, useState } from 'react'
import { Formik, Form, useField } from 'formik'
import * as Yup from 'yup'
import { BrowserRouter, Routes, Route, useNavigate, Link } from 'react-router-dom'
import { DesignSystemPage } from './page/design-system'
import { BorderlessButton } from './view/button'

interface KDTextInputProps {
  label: string
  name: string
  type: string
  placeholder?: string
  id?: string
}
const KDTextInput = ({ ...props }: KDTextInputProps): ReactElement => {
  // useField() returns [formik.getFieldProps(), formik.getFieldMeta()]
  // which we can spread on <input> and alse replace ErrorMessage entirely.
  const [field, meta] = useField(props)
  const errorInput = meta.touched && meta.error !== undefined

  return (
    <div className='w-full flex flex-row'>
      <label className='text-white' htmlFor={props.id ?? props.name}>{props.label}</label>
      <input className="text-black text-input grow text-right" {...field} {...props} />
      {errorInput
        ? <div className="tooltip tooltip-open tooltip-error absolute" data-tip={meta.error}></div>
        : null}
    </div>
  )
}

interface KDCheckboxProps {
  name: string
  children: string
  id?: string
}
const KDCheckbox = ({ children, ...props }: KDCheckboxProps): ReactElement => {
  // useField() returns [formik.getFieldProps(), formik.getFieldMeta()]
  // which we can spread on <input> and alse replace ErrorMessage entirely.
  const [field, meta] = useField({ ...props, type: 'checkbox' })
  return (
    <div className='m-auto'>
      <label className="checkbox">
        <input {...field} {...props} type="checkbox" />
        {children}
      </label>
      {meta.touched && meta.error !== undefined
        ? (
        <div className="error">{meta.error}</div>
          )
        : null}
    </div>
  )
}

interface KDSelectProps {
  label: string
  name: string
  id?: string
  children: ReactElement[]
}
const KDSelect = ({ label, ...props }: KDSelectProps): ReactElement => {
  // useField() returns [formik.getFieldProps(), formik.getFieldMeta()]
  // which we can spread on <input> and alse replace ErrorMessage entirely.
  const [field, meta] = useField(props)
  return (
    <div className='m-auto'>
      <label htmlFor={props.id ?? props.name}>{label}</label>
      <select {...field} {...props} />
      {meta.touched && meta.error !== undefined
        ? (
        <div>{meta.error}</div>
          )
        : null}
    </div>
  )
}

const SignUpForm = (): ReactElement => {
  return (
    <div className='text-center'>
      <h1>Sign Up</h1>
      <Formik
        initialValues={{
          name: '',
          email: '',
          password: '',
          acceptedTerms: false, // added for our checkbox
          jobType: '' // added for our select
        }}
        validationSchema={Yup.object({
          name: Yup.string()
            .max(15, 'Must be 15 characters or less')
            .required('Required'),
          email: Yup.string()
            .email('Invalid email`')
            .required('Required'),
          password: Yup.string()
            .required('Required'),
          acceptedTerms: Yup.boolean()
            .required('Required')
            .oneOf([true], 'You must accept the terms and conditions.'),
          jobType: Yup.string()
            // specify the set of valid values for job type
            // @see http://bit.ly/yup-mixed-oneOf
            .oneOf(
              ['designer', 'development', 'product', 'other'],
              'Invalid Job Type'
            )
            .required('Required')
        })}
        onSubmit={async (values, { setSubmitting }) => {
          await new Promise((resolve, reject) => setTimeout(resolve, 1000))
          console.log('kodok', values)
          setSubmitting(false)
        }}
      >
        <Form className='flex flex-col'>
          <KDTextInput
            label="Name"
            name="name"
            type="text"
            placeholder="Kay Dee"
          />
          <KDTextInput
            label="Email"
            name="email"
            type="email"
            placeholder="local-kd@yopmail.com"
          />
          <KDTextInput
            label="Password"
            name="password"
            type="password"
          />
          <KDSelect label="Job Type" name="jobType">
            <option value="">Select a job type</option>
            <option value="designer">Designer</option>
            <option value="development">Developer</option>
            <option value="product">Product Manager</option>
            <option value="other">Other</option>
          </KDSelect>
          <KDCheckbox name="acceptedTerms">
            I accept the terms and conditions
          </KDCheckbox>

          <button type="submit">Submit</button>
        </Form>
      </Formik>
    </div>
  )
}

interface SignInFormProps {
  onSubmit: (user: UserModel) => void
}
const SignInForm = (props: SignInFormProps): ReactElement => {
  const navigate = useNavigate()
  return (
    <div className='text-center'>
      <h1>Sign In</h1>
      <Formik
        initialValues={{
          email: '',
          password: ''
        }}
        validationSchema={Yup.object({
          email: Yup.string()
            .email('Invalid email addresss`')
            .required('Required'),
          password: Yup.string()
            .required('Required')
          // add more validation with yup-password library later ~kodok
          // Yup.string().required('Required').min(8, 'password must contain 8 or more characters with at least one of each: uppercase, lowercase, number and special').minLowercase(1, 'password must contain at least 1 lower case letter').minUppercase(1, 'password must contain at least 1 upper case letter').minNumbers(1, 'password must contain at least 1 number').minSymbols(1, 'password must contain at least 1 special character');
        })}
        onSubmit={async (values, { setSubmitting }) => {
          // use preventDefault? ~kodok
          await new Promise((resolve, reject) => setTimeout(resolve, 1000))
          const backendResponse = { name: 'Kay Dee', ...values }
          console.log('kodok', values)
          props.onSubmit(backendResponse)
          setSubmitting(false)
          navigate('/home', { replace: true })
        }}
      >
        <Form className='flex flex-col'>
          <KDTextInput
            label="Email"
            name="email"
            type="email"
            placeholder="local-kd@yopmail.com"
          />
          <KDTextInput
            label="Password"
            name="password"
            type="password"
          />

          <button type="submit">Submit</button>
        </Form>
      </Formik>
    </div>
  )
}

interface UserModel {
  name: string
  email: string
}
interface HomeProps {
  user?: UserModel
}
const Home = ({ user }: HomeProps): ReactElement => {
  const greeting = user !== undefined ? `Hello again ${user.name}!` : 'Hello stranger!'
  return <h1 className='text-center'>{greeting}</h1>
}

interface NavProps {
  user?: UserModel
  setUser: (user?: UserModel) => void
}
const Nav = (props: NavProps): ReactElement => {
  const logout = (): void => {
    props.setUser(undefined)
  }
  const links = props.user !== undefined
    ? (<>
        <Link to="/home" className="grow" onClick={logout}>Sign out</Link>
      </>)
    : (<>
        <Link to="/signup" className='grow'><BorderlessButton>Sign up</BorderlessButton></Link>
        <Link to="/signin" className='grow'><BorderlessButton>Sign in</BorderlessButton></Link>
      </>)

  return (<div className='flex flex-row text-center'>
    {links}
  </div>)
}

function App (): ReactElement {
  const [user, setUser] = useState<UserModel | undefined>(undefined)

  return (
    <div className="App">
      <BrowserRouter>
        <Nav user={user} setUser={setUser}/>
        <Routes>
          <Route path="/" element={<Home user={user}/>}/>
          <Route path="/home" element={<Home user={user}/>}/>
          <Route path="/signin" element={<SignInForm onSubmit={setUser}/>}/>
          <Route path="/signup" element={<SignUpForm/>}/>
          <Route path="/design-system" element={<DesignSystemPage/>}/>
        </Routes>
      </BrowserRouter>
    </div>
  )
}

export default App
