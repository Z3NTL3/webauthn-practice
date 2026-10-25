import { Outlet } from "react-router";
import Logo from "@/assets/svg/logo.svg?react";

export default function Navbar() {
    return (
        <>
            <div className="fixed flex items-center top-0 left-0  w-full min-h-15">
                <div onClick={() => window.location.href = "/"} className="cursor-pointer flex items-center justify-center gap-x-2 ml-4 mt-4  md:ml-8 md:mt-8">
                    <div className="bg-white p-3 rounded-xl">
                        <Logo/> 
                    </div>
                    <h2 className="font-medium text-lg">Passkey Demo</h2>
                </div>

                <div className="md:flex grow justify-end mr-8 mt-8 hidden">
                    <div className="flex items-center justify-center py-1 gap-x-2 px-4 border rounded-full text-md bg-surface-primary text-muted border-tint">
                        <div className="bg-green-800 rounded-full w-2 h-2"></div>
                        WebAuthn - Demo
                    </div>
                </div>
            </div>
            <Outlet />
        </>
    )
}