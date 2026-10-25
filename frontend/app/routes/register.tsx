import { createRef } from "react";
import type { Route } from "./+types/home";
import { HiUserAdd } from "react-icons/hi";
import { FaFingerprint } from "react-icons/fa6";
import { TbShieldBolt } from "react-icons/tb"; 
import { toast } from "~/components/ui/toast"

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Register - WebAuthn Practice Demo" }
  ];
}

export default function Register() {
    const inputRef = createRef<HTMLInputElement>()  

    const handleRegister = async() => {
        try {
            const email = inputRef.current?.value;
            if(!email) 
                return toast.add({
                    description: "Geen geldige email adres",
                    type: "error"
                });
            
            let req = await fetch(`${import.meta.env.VITE_BACKEND_URL}/register/create`, {
                method: "POST",
                credentials: "include",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({email})
            });
            if(!req.ok) {
                toast.add({
                    description: "Er ging iets mis bekijk de console",
                    type: "error"
                })
                throw new Error(`failed creating options: ${req.statusText} ${await req.text()}`)
            }   
    
            let body = await req.json()
            const publicKey = PublicKeyCredential.parseCreationOptionsFromJSON(body.data.publicKey);
            let credential = await navigator.credentials.create({
                publicKey,
            })

            console.log("credential", credential)

            if(!credential) {
                 toast.add({
                    description: "Er ging iets mis bekijk de console",
                    type: "error"
                })
                throw new Error("failed retrieving credential response")
            }

            req = await fetch(`${import.meta.env.VITE_BACKEND_URL}/register/finish`, {
                method: "POST",
                credentials: "include",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({email, credential})
            });
            if(!req.ok) {
                 toast.add({
                    description: "Er ging iets mis bekijk de console",
                    type: "error"
                })
                throw new Error(`failed finalizing registiration handshake: ${req.statusText} ${await req.text()}`)
            }   

            console.log("registiration finalization req: ",req)
            toast.add({
                description: "Registratie met success, ga verder naar het inlog scherm",
                type: "success"
            })
        } catch (err) {
             toast.add({
                description: "Er ging iets mis bekijk de console",
                type: "error"
            })
            console.error(err)
        }
    }

    return <div className="flex w-screen h-screen justify-center items-center px-2 md:px-0">
        <div className="flex flex-col border border-tint bg-surface-primary text-white p-8 w-full h-fit md:w-110 md:min-h-100 rounded-2xl shadow-[4px_4px_40px_4px_rgba(0,0,0,0.1)]">
        <div className="bg-accent w-fit h-fit p-4 rounded-lg">
            <HiUserAdd size={24} color="#A5B4FC" />
        </div>

            <div className="mt-5"/>
            <h1 className="font-bold text-2xl">Registratie</h1>
            <p className="text-muted">Registreer hier voor een veilige sleutel.</p>

            <div className="mt-10"/>
            <h3 className="font-semibold">Email</h3>
            <input ref={inputRef} className="border border-tint p-2 rounded-lg outline-[#2A55C9]" type="email" placeholder="hello@world.com" required/>

            <button onClick={handleRegister} className="flex justify-center items-center bg-text py-2 gap-x-2 px-5 text-black rounded-lg mt-4 cursor-pointer">
             <FaFingerprint color="black" size={20} /> Maak sleutel aan
            </button>

            <span className="flex text-[14px] items-center justify-center gap-x-2 text-muted mt-5 self-start">
                <TbShieldBolt size={20}/> Straks log je veilig in met een passkey.
            </span>
            <span className="text-center text-muted mt-8">Heb je al een account <a href="/" className="text-blue-500 underline cursor-pointer">klik dan hier</a> om in te loggen.</span>
        </div>
    </div>
}
