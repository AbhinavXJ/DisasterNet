import { AlertTriangle, Send } from 'lucide-react'
import React, { useEffect, useState } from 'react'

function App() {
  const [messages, setMessages] = useState([])
  const [message, setMessage] = useState('')

  useEffect(()=>{
    const interval = setInterval(()=>{
      fetch('http://localhost:3001/messages')
      .then(res=>res.json())
      .then(data => {
        console.log(data) 
        setMessages(data)
      }
      )
    },2000)
    return ()=> clearInterval(interval)
  },[])



  const sendMessage = async () => {
    if (message.trim() === '') return

    await fetch('http://localhost:3001/send',{
      method:'POST',
      headers:{'Content-Type':'application/json'},
      body: JSON.stringify({message:message}),
    })
    setMessage('')
  }
  return (
<div className='min-h-screen w-full flex flex-col justify-center items-center bg-gradient-to-tr from-violet-950 to-violet-500'>

      <div className='flex flex-col h-[80vh] sm:w-[80vw] lg:w-[60vw]  border-4 border-black rounded-3xl overflow-hidden my-8'>

          <div className='flex h-[10%] bg-red-700'>
                  <div className='my-auto ml-4'>
                  <AlertTriangle className="text-yellow-400 w-12 h-12" />
                </div>
                <div className='text-white ml-4 my-2'>
                   <p className='text-3xl font-bold'>DisasterNet  (No Internet Required)</p>
                   <p className='text-md '>Emergency Communication Network</p>
                </div>
          </div>

           {/* MESSAGES SECTION */}
        <div className='h-[80%] bg-slate-700 px-6 py-4 overflow-y-auto space-y-2'>
        {Array.isArray(messages) ? (
  messages.map((msg, idx) => (
    <div key={idx} className='text-white bg-slate-800 px-4 py-2 rounded-md'>
      {msg}
    </div>
  ))
) : (
  <div className="text-red-400">No messages available</div>
)}
        </div>

 {/* INPUT SECTION */}
        <div className='h-[10%] bg-slate-800 flex items-center px-6'>
          <input
            type='text'
            placeholder='Type your message here...'
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            className='w-full p-2 rounded-lg bg-slate-600 text-white'
          />
          <button onClick={sendMessage} className='ml-4'>
            <Send className='text-gray-300 w-8 h-8' />
          </button>
        </div>



      </div>
      <div className="mt-6 px-4 w-full flex justify-center my-8">
  <div className="max-w-2xl w-full bg-white/5 backdrop-blur-lg rounded-2xl p-6 shadow-md text-center">
    <p className="text-xl md:text-2xl font-bold mb-2">
      Created By -{" "}
      <a
        href="https://github.com/AbhinavXJ"
        target="_blank"
        rel="noopener noreferrer"
        className="text-blue-400 hover:underline"
        aria-label="AbhinavXJ"
      >
        Abhinav Jha
      </a>
    </p>

    <p className="text-base md:text-lg text-gray-300 mb-4">
      Frontend part of the <strong>DisasterNet</strong> project — made to
      communicate in places <span className="font-semibold text-white">without the NEED OF INTERNET</span>!
    </p>

    <a
      href="https://github.com/AbhinavXJ/DisasterNet"
      target="_blank"
      rel="noopener noreferrer"
      className="inline-block text-blue-400 font-medium hover:underline"
      aria-label="DisasterNet GitHub Repository"
    >
      View full project (frontend + backend) on GitHub →
    </a>
  </div>
</div>


    </div>
  )
}

export default App