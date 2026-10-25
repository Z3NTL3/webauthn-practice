import type { Route } from "./+types/home";
import { IoKey as  Key} from "react-icons/io5";
import { LuKeyRound } from "react-icons/lu";
import { createRef, useState } from "react";
import { TbShieldCheckFilled } from "react-icons/tb";
import { toast } from "~/components/ui/toast"

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Login - WebAuthn Practice Demo" }
  ];
}

type User = {id: string, pubKey: string}
export default function Home() {
  const inputRef = createRef<HTMLInputElement>()  
  const [user, setUser] = useState<User>();

  const handleLogin = async() => {
      try {
          const email = inputRef.current?.value;
          if(!email)
            return toast.add({
                description: "Vul een geldige email adres in",
                type: "error"
            })
          
          let req = await fetch(`${import.meta.env.VITE_BACKEND_URL}/login/begin`, {
              method: "POST",
              credentials: "include",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify(email)
          });
          if(!req.ok) {
            toast.add({
                description: "Er ging iets mis bekijk de console",
                type: "error"
            })
            throw new Error(`failed creating options: ${req.statusText} ${await req.text()}`)
          }   
  
          let body = await req.json()
          const publicKey = PublicKeyCredential.parseRequestOptionsFromJSON(body.data.publicKey);
          let assertion = await navigator.credentials.get({
              publicKey,
          })

          console.log("assertion", assertion)

          req = await fetch(`${import.meta.env.VITE_BACKEND_URL}/login/finalize`, {
              method: "POST",
              credentials: "include",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify({email, assertion})
          });
          if(!req.ok) {
            toast.add({
                description: "Er ging iets mis bekijk de console",
                type: "error"
            })
            throw new Error(`failed creating options: ${req.statusText} ${await req.text()}`)
          }   

          body = await req.json()
          console.log("credentials valid; signed in!")

          toast.add({
            description: "Je bent succesvol ingelogd",
            type: "success"
          })
          setUser({
            id: body.data.id,
            pubKey: body.data.pubKey
          })
      } catch (err) {
        toast.add({
          description: "Er ging iets mis bekijk de console",
          type: "error"
        })
        console.error(err)
      }
  }

  const handleSignOff = async() => {
    setUser(() => undefined)
  }

  return <div className="flex w-screen h-screen justify-center items-center px-2 md:px-0">
      <div className={`flex flex-col ${user ? "items-center" : null} border border-tint bg-surface-primary text-white p-8 w-full h-fit md:w-110 md:min-h-100 rounded-2xl shadow-[4px_4px_40px_4px_rgba(0,0,0,0.1)]`}>
        <div className={`w-fit h-fit p-4 rounded-lg ${user ? "bg-success" : "bg-accent"}`}>
            {user ? <TbShieldCheckFilled size={24} color="#6EE7B7"  /> : <Key size={24} color="#A5B4FC" />}
        </div>

        <div className="mt-5"/>
        <h1 className="font-bold text-2xl">{user ? "Je bent ingelogd" : "Log in"}</h1>
        <p className="text-muted">{user ? "Je passkey is succesvol geverifieerd." : "Log hier veilig in met een passkey"}</p>

        {!user ? 
          <>
            <div className="mt-10"/>
            <h3 className="font-semibold">Email</h3>
            <input ref={inputRef} className="border border-tint p-2 rounded-lg outline-[#2A55C9]" type="email" placeholder="hello@world.com" required/>

            <button onClick={handleLogin} className="flex justify-center items-center bg-text py-2 gap-x-2 px-5 text-black rounded-lg mt-4 cursor-pointer">
              <LuKeyRound size={24}  /> Log in met een sleutel
            </button>

            <span className="text-center mt-5 text-muted">Heb je nog geen account <a href="/register" className="text-blue-500 underline cursor-pointer">klik dan hier</a> om te registreren.</span>
          </> : 
          <>
            <div className="bg-[#1F232B] rounded-lg h-fit w-full mt-10 py-5 px-5">
              <div className="flex items-center justify-center gap-x-2">
                <h2 className="font-semibold text-muted">ID</h2>
                <span className="text-xs text-[#E6E8EE]">{user.id}</span>
              </div>
            </div>

            <button onClick={handleSignOff} className="flex justify-center items-center py-2 gap-x-2 px-25 mt-10 text-white border border-tint rounded-lg cursor-pointer">
              Log uit
            </button>
          </>
        }
        </div>
  </div>
}
