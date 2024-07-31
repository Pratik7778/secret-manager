import React, { useState, useEffect } from 'react';
import { Button, Modal, Form } from 'react-bootstrap';
import Cookies from 'js-cookie';

const CreateSecretPage = ({refreshFlag, setRefreshFlag}) => {
    const [showModal, setShowModal] = useState(false);
    const [key, setKey] = useState('');
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
            const response = await fetch(API_URL + `/api/v1/secrets/create`, {
                method: "POST",
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    "key": key,
                    "value":value
                })
            })

            const result = await response.json()

            if(response.status === 200) {
                setMessage({type: "success", message: result.message})
                setRefreshFlag(!refreshFlag)
                setKey("")
                setValue("")

                console.log(message)

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

        console.log(`Key: ${key}, Value: ${value}`);
        // Close the modal
        handleClose();
    };

    useEffect(() => {
        if(message) {
            setTimeout(() => {
                setMessage("")
            }, 5000)
        }
    }, [message])

    return (
        <div className="container mt-5">
            <div className="d-flex justify-content-end mb-3">
                <Button variant="primary" onClick={handleShow}>
                    Create Secret
                </Button>
            </div>

        <Modal show={showModal} onHide={handleClose}>
            <Modal.Header closeButton>
            <Modal.Title>Enter Key and Value</Modal.Title>
            </Modal.Header>
            <Modal.Body>
            <Form onSubmit={handleSubmit}>
                <Form.Group controlId="formKey">
                <Form.Label>Key</Form.Label>
                <Form.Control
                    type="text"
                    placeholder="Enter key"
                    value={key}
                    onChange={(e) => setKey(e.target.value)}
                />
                </Form.Group>

                <Form.Group controlId="formValue">
                <Form.Label>Value</Form.Label>
                <Form.Control
                    type="text"
                    placeholder="Enter value"
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
            <Button variant="secondary" onClick={handleClose}>
                Close
            </Button>
            </Modal.Footer>
        </Modal>

        {message && (
            <div className={`alert ${message.type === 'error' ? 'alert-danger' : 'alert-success'} mt-3`}>
                {message.message}
            </div>
        )}
        </div>
    );
};

export default CreateSecretPage;
