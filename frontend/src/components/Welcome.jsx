import React from "react";
import { useState, useEffect } from "react";
import { useNavigate } from 'react-router-dom';


const Homepage = () => {
    // const [message, setMessage] = useState("")  
    const [apiURL, setApiURL] = useState("http://localhost:8080")
    const navigate = useNavigate()

    // const fetchWelcomeMsg = async () => {
    //     try {
    //         const API_URL = await import.meta.env.VITE_API_URL
    //         // const API_URL = process.env.API_URL
    //         const response = await fetch(API_URL + "/api/v1", {
    //             method: "GET",
    //             headers: {
    //                 'Content-Type': 'application/json',
    //             },
    //         })
            
    //         if(!response.ok) {
    //             throw new Error("Network response was not ok")
    //         }

    //         // console.log("Response: ", response.data)
    //         const result = await response.json()
    //         setMessage(result)

    //         console.log("Result: ", result)

    //     } catch(err) {
    //         console.log(err)
    //     }
    // }

    // useEffect(() => {
    //     fetchWelcomeMsg()
    // }, [])

    const handleInputChange = (e) => {
        setApiURL(e.target.value)
    }

    const handleSubmit = (e) => {
        e.preventDefault()
        
        localStorage.setItem("API_URL", apiURL)
        navigate("/login")
    }

    return (
        <>
        <div className="col-md-12">
        <div className="card card-container">
            <div className="container mt-5">
                <h2>Enter API URL</h2>
                <form onSubmit={handleSubmit}>
                    
                    <div className="mb-3">
                        {/* <label htmlFor="username" className="form-label">Username:</label> */}
                        <input
                            type="text"
                            id="url"
                            name="url"
                            value={apiURL}
                            onChange={handleInputChange}
                            className="form-control"
                            // required
                        />
                    </div>

                    <button type="submit" className="btn btn-primary">Confirm</button>
                </form>
            </div>
        </div>
        </div>
            {/* {message} */}
        </>
    )
}

export default Homepage;