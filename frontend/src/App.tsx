import { useState } from "react";

import "./App.css";
import Quote from "./components/Quote";

function App() {
  const [count, setCount] = useState(0);

  return (
    <>
      <Quote></Quote>
    </>
  );
}

export default App;
