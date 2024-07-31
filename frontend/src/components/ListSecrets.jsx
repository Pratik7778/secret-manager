import React, {useState, useEffect} from "react";
import ViewSecretsPage from "./ViewSecrets";
import DeleteSecretPage from "./DeleteSecrets";
import CreateSecretPage from "./CreateSecret";
import UpdateSecretPage from "./UpdateSecret";

import { useNavigate } from 'react-router-dom';
import LogoutPage from "./Logout";
import Cookies from 'js-cookie';


const ListSecretsPage = () => {
    const [secrets, setSecrets] = useState([])
    const [query, setQuery] = useState("")
    const [page, setPage] = useState(1)
    const [pageSize, setPageSize] = useState(10)
    const [totalRecords, setTotalRecords] = useState(0)
    const [message, setMessage] = useState(null)
    const [refreshFlag, setRefreshFlag] = useState(false)

    const navigate = useNavigate()
    
    useEffect(() => {
        const token = Cookies.get("token")
        if(!token) {
            navigate("/login")
        }            
    }, [navigate])

    const handlePageChange = (newPage) => {
        setPage(newPage)
    }

    const handlePageSizeChange = (e) => {
        setPageSize(Number(e.target.value))
    }

    const fetchSecrets = async () => {
        const token = Cookies.get("token")
        if(!token) {
            setMessage({type: "error", message: "Unauthorized: No token found."})
            return
        }

        try {
            // const API_URL = await import.meta.env.VITE_API_URL
            const API_URL = localStorage.getItem("API_URL") || "http://localhost:8080"
            const response = await fetch(API_URL + `/api/v1/secrets?q=${query}&page=${page}&page_size=${pageSize}`, {
                method: "GET",
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
            })
            
            const result = await response.json()

            if(response.status === 200) {
                setSecrets(result.secrets)
                setTotalRecords(result.total)

                console.log(result)

            } else if(response.status === 401) {
                setMessage({type: "error", message: result})
                console.log(result)

            } else if(response.status === 404) {
                setMessage({type: "error", message: result.error})
                console.log(result)

            } else {
                throw new Error("Network response was not ok")
            }

        } catch(err) {
            console.log(err)
        }
    }

    useEffect(() => {
        fetchSecrets()
    }, [query, page, pageSize, refreshFlag])

    useEffect(() => {
        if(message) {
            setTimeout(() => {
                setMessage("")
            }, 5000)
        }
    }, [message])

    return (
        <>
            <div className="container mt-5">
                <h2>Secrets List</h2>
                <LogoutPage />
                <CreateSecretPage refreshFlag={refreshFlag} setRefreshFlag={setRefreshFlag}/>
                <div className="mb-3">
                    <label htmlFor="search" className="form-label">Search:</label>
                    <input 
                        type="text"
                        id="search"
                        value={query}
                        onChange={(e) => setQuery(e.target.value.trim())}
                        className="form-control"
                    />

                    {message && (
                        <div className={`alert ${message.type === "error" ? "alert-danger" : "alert-success"} mt-3`}>
                            {message.message}
                        </div>
                    )}

                    <table className="table">
                        <thead>
                            <tr>
                                <th>Secret Name</th>
                                <th>Actions</th>
                            </tr>
                        </thead>

                        <tbody>
                            {secrets.map((secret, index) => (
                                <tr key={index}>
                                    <td>{secret}</td>
                                    
                                    <ViewSecretsPage secretName={secret} />
                                    <UpdateSecretPage secretName={secret} />
                                    <DeleteSecretPage secretName={secret} refreshFlag={refreshFlag} setRefreshFlag={setRefreshFlag}/>
                                
                                </tr>
                            ))}
                        </tbody>
                    </table>

                    <div className="d-flex justify-content-between">
                        <div>
                            <label htmlFor="pageSize" className="form-label">Page Size:</label>
                            <select id="pageSize" value={pageSize} onChange={handlePageSizeChange} className="form-select">
                                <option value={5}>5</option>
                                <option value={10}>10</option>
                                <option value={20}>20</option>
                            </select>   
                        </div>
                        
                        <div>
                            <button
                                className="btn btn-primary me-2"
                                disabled={page <= 1}
                                onClick={() => handlePageChange(page-1)}
                            >Previous</button>
                            <span>Page {page}</span>
                            <button
                                className="btn btn-primary ms-2"
                                disabled={page * pageSize >= totalRecords}
                                onClick={() => handlePageChange(page+1)}
                            >Next</button>
                        </div>
                    </div>

                    <div className="mt-3">
                        <span>Total Records: {totalRecords}</span>
                    </div>

                </div>
            </div> 
        </>
    )
}

export default ListSecretsPage;