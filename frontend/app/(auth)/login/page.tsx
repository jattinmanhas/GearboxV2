import Link from 'next/link'
import React from 'react'
import LoginForm from './components/LoginForm'

const LoginPage = () => {
  return (
    <>
    <div className="flex flex-col items-center justify-center m-2 md:m-10">
      <h3 className="text-xl font-semibold text-neutral-400">GearBox</h3>
      <h1 className="text-3xl font-bold">Welcome Back</h1>
      <LoginForm />
      <Link className="py-3" href="/register">
        Dont have an account?
      </Link>
    </div>
  </>

  )
}

export default LoginPage