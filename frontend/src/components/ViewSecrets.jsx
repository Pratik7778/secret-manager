import React, {useState} from "react";
import {Modal} from 'react-bootstrap';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faEye } from '@fortawesome/free-solid-svg-icons';
import Cookies from 'js-cookie';


const ViewSecretsPage = ({secretName}) => {
    const [secrets, setSecrets] = useState({})
    const [message, setMessage] = useState(null)

    const [showModal, setShowModal] = useState(false);
    const handleShow = () => {
        setShowModal(true)
        fetchSecretContent()
    }
    const handleClose = () => setShowModal(false);

    const fetchSecretContent = async () => {
        const token = Cookies.get("token")
        if(!token) {
            setMessage({type: "error", message: "Unauthorized: No token found."})
            return
        }

        try {
            // const API_URL = await import.meta.env.VITE_API_URL
            const API_URL = localStorage.getItem("API_URL") || "http://localhost:8080"
            const response = await fetch(API_URL + `/api/v1/secrets/${secretName}`, {
                method: "GET",
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
            })

            const result = await response.json()

            if(response.status === 200) {
                setSecrets(result)
                console.log(result)

            } else if(response.status === 400) {
                setMessage({type: "error", message: result.error})
                console.log(result)

            } else if(response.status === 401) {
                setMessage({type: "error", message: result}) 
                console.log(result)

            } else {
                throw new Error("Network response was not ok")
            }

        } catch(err) {
            console.log(err)
        }
    }

    return (
        <>
            <td>
                <FontAwesomeIcon
                    icon={faEye}
                    size="1x"
                    onClick={handleShow}
                    style={{cursor: "pointer", marginRight: "-100px"}}
                />

                <Modal show={showModal} onHide={handleClose}>
                    <Modal.Header closeButton>
                    <Modal.Title>Secret Value</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        {secrets && (
                            <p>{secrets[secretName]}</p>
                        )}
                    </Modal.Body>
                    <Modal.Footer>
                    </Modal.Footer>
                </Modal>
            </td>
        </>
    )
}

export default ViewSecretsPage;
