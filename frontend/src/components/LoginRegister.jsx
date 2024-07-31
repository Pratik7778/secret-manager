import React from "react";
import { useState, useEffect } from "react";
import { useNavigate } from 'react-router-dom';
import Cookies from 'js-cookie';

const LoginRegisterPage = () => {
    const [formData, setFormData] = useState({
        user: "",
        password: ""
    })
    const [message, setMessage] = useState("")
    const [action, setAction] = useState("")
    
    const navigate = useNavigate()

    useEffect(() => {
        const token = Cookies.get("token")
        if(token) {
            navigate("/secrets")
        }            
    }, [navigate])

    const handleInputChange = (e) => {
        const {name, value} = e.target
       
        setFormData({
            ...formData,
            [name]: value
        })
    }

    const handleSubmit = async (e) => {   
        e.preventDefault()

        try {
            // const API_URL = await import.meta.env.VITE_API_URL
            const API_URL = localStorage.getItem("API_URL") || "http://localhost:8080"
            const url = action === "register" ? API_URL + "/api/v1/register" : API_URL + "/api/v1/login"

            const response = await fetch(url, {
                method: "POST",
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            })
            
            const result = await response.json()

            if(response.status === 200) {
                if(action === "register") {
                    setMessage({type: "success", message: result.message})
                } else if (action === "login" && result.token) {
                    Cookies.set('token', result.token, { expires: 1 })
                    console.log(Cookies.get('token'))

                    // localStorage.setItem("token", result.token)
                    // console.log(localStorage.getItem("token"))

                    //Go to list secrets page
                    navigate("/secrets")
                }

                console.log(result)
                setFormData({
                    user: "",
                    password: ""
                })
            } else if(response.status === 400) {
                setMessage({type: "error", message: result.error})
                console.log(result)
            } else {
                throw new Error("Network response was not ok")
            }

            // if(!response.ok) {
            //     throw new Error("Network response was not ok")
            // }

        } catch(err) {
            console.log(err)
        }
    }

    useEffect(() => {
        if(message) {
            setTimeout(() => {
                setMessage("")
            }, 5000)
        }
    }, [message])

    return (
        <>
        <div className="col-md-12">
            <div className="card card-container">
                <div className="container mt-5">
                    <h2>Secret Manager</h2>
                    <form onSubmit={handleSubmit}>
                        
                        <div className="mb-3">
                            <label htmlFor="username" className="form-label">Username:</label>
                            <input
                                type="text"
                                id="user"
                                name="user"
                                value={formData.user}
                                onChange={handleInputChange}
                                className="form-control"
                                required
                            />
                        </div>
                        
                        <div className="mb-3">
                            <label htmlFor="password" className="form-label">Password:</label>
                            <input
                                type="password"
                                id="password"
                                name="password"
                                value={formData.password}
                                onChange={handleInputChange}
                                className="form-control"
                                required
                            />
                        </div>
                        <button type="submit" className="btn btn-primary" onClick={() => setAction("login")}>Login</button>
                        <button type="submit" className="btn btn-secondary ms-2" onClick={() => setAction("register")}>Register</button>

                    </form>
                </div>

                {message && (
                    <div className={`alert ${message.type === 'error' ? 'alert-danger' : 'alert-success'} mt-3`}>
                        {message.message}
                    </div>
                )}

            </div>
        </div>
        </>
    )
}

export default LoginRegisterPage;