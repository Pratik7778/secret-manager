import React, { useEffect, useState } from 'react';
import { Button, Modal, Form } from 'react-bootstrap';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPen } from '@fortawesome/free-solid-svg-icons';
import Cookies from 'js-cookie';

const UpdateSecretPage = ({secretName, refreshFlag, setRefreshFlag}) => {
    const [showModal, setShowModal] = useState(false);
    const [value, setValue] = useState('');
    const [message, setMessage] = useState(null)


    const handleShow = () => setShowModal(true);
    const handleClose = () => setShowModal(false);

    const handleSubmit = async (e) => {
        e.preventDefault();

        const token = Cookies.get("token")
        if(!token) {
            setMessage({type: "error", message: "Unauthorized: No token found."})
            return
        }

        try {
            // const API_URL = await import.meta.env.VITE_API_URL
            const API_URL = localStorage.getItem("API_URL") || "http://localhost:8080"
            const response = await fetch(API_URL + `/api/v1/secrets/${secretName}`, {
                method: "PUT",
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    "value":value
                })
            })

            const result = await response.json()

            if(response.status === 200) {
                setMessage({type: "success", message: result.message})
                // setValue("")

                console.log(result)                

            } else if(response.status === 400) {
                setMessage({type: "error", message: result.error})
                console.log(result)

            }else if(response.status === 401) {
                setMessage({type: "error", message: result}) 
                console.log(result)

            } else {
                throw new Error("Network response was not ok")
            }

        } catch(err) {
            console.log(err)
        }

        console.log(`Value: ${value}`);
        // Close the modal
        handleClose();
    };

    const fetchSecretContent = async () => {
        const token = Cookies.get("token")
        if(!token) {
            setMessage({type: "error", message: "Unauthorized: No token found."})
            return
        }

        try {
            const API_URL = await import.meta.env.VITE_API_URL
            // const API_URL = process.env.API_URL
            const response = await fetch(API_URL + `/api/v1/secrets/${secretName}`, {
                method: "GET",
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
            })

            const result = await response.json()

            if(response.status === 200) {
                setValue(result[secretName])
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

    useEffect(() => {
        fetchSecretContent()
    }, [])

    return (
        <>
        <td>
            <FontAwesomeIcon
                icon={faPen}
                size="1x"
                onClick={handleShow}
                style={{cursor: "pointer", marginLeft: "-400px"}}
            />

            <Modal show={showModal} onHide={handleClose}>
                <Modal.Header closeButton>
                <Modal.Title>Enter Value</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                <Form onSubmit={handleSubmit}>
                    <Form.Group controlId="formValue">
                    <Form.Label>Value</Form.Label>
                    <Form.Control
                        type="text"
                        placeholder={value}
                        value={value}
                        onChange={(e) => setValue(e.target.value)}
                    />
                    </Form.Group>

                    <Button variant="primary" type="submit" style={{marginTop: "20px"}}>
                    Submit
                    </Button>
                </Form>
                </Modal.Body>
                <Modal.Footer>
                </Modal.Footer>
            </Modal>
        </td>
        </>
    );
};

export default UpdateSecretPage;
