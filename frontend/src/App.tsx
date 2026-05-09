import { useEffect, useState } from "react";
import { Quit } from "../wailsjs/runtime/runtime";

import "./App.css";

import { GetConnections, GetLinkData } from "../wailsjs/go/main/App";
import { main } from "../wailsjs/go/models";

function App() {
const [connections, setConnections] = useState<main.ConnectDTO[]>([]);
const [selectedIndex, setSelectedIndex] = useState<number | null>(null);
const [results, setResults] = useState<main.LinkDataDTO | null>(null);
``

  useEffect(() => {
    GetConnections().then(setConnections);
  }, []);

  const selected =
    selectedIndex !== null ? connections[selectedIndex] : null;

  const handleSelect = (index: number) => {
    setSelectedIndex(index);
    setResults(null);
  };

  const handleGetLinkData = async () => {
    if (!selected) return;
    const data = await GetLinkData((selected as any).id);
    setResults(data);
  };

  return (
    <div id="App">
      <div className="panel-container">

        {/* Selection */}
        <div className="boxed">
          <h3 className="box-title">Selection</h3>

          <div className="selection-grid">
            <label>Network Connection:</label>
            <select
              value={selectedIndex ?? ""}
              onChange={e =>
                handleSelect(Number(e.target.value))
              }
            >
              <option value="">Please select...</option>
              {connections.map((c: any, index: number) => (
                <option key={index} value={index}>
                  {c.name}
                </option>
              ))}
            </select>

            <label>Network Card:</label>
            <span>{(selected as any)?.network_card ?? "-"}</span>

            <label>MAC Address:</label>
            <span>{(selected as any)?.mac_addr ?? "-"}</span>

            <label>IP Address:</label>
            <span>{(selected as any)?.ip_addr || "-"}</span>
          </div>

          <div className="button-row">
            <button
              onClick={handleGetLinkData}
              disabled={!selected}
            >
              Get Link Data
            </button>
            <button>Save Link Data</button>
            <button>Help</button>
            <button onClick={Quit}>Quit</button>
          </div>
        </div>

        {/* Results */}
        <div className="boxed results-box">
          <h3 className="box-title">Results</h3>

          <div className="results-grid">
            <label>Switch Name:</label>
            <span>{(results as any)?.switch_name ?? "-"}</span>

            <label>Switch IP Address:</label>
            <span>{(results as any)?.switch_ip ?? "-"}</span>

            <label>Port Identifier:</label>
            <span>{(results as any)?.port_id ?? "-"}</span>

            <label>Switch Model:</label>
            <span>{(results as any)?.switch_model ?? "-"}</span>

            <label>VLAN Identifier:</label>
            <span>{(results as any)?.vlan_id ?? "-"}</span>

            <label>Port Duplex:</label>
            <span>{(results as any)?.duplex_option ?? "-"}</span>

            <label>VPT Management Domain:</label>
            <span>{(results as any)?.vpt_mgmt_domain ?? "-"}</span>

            <label>Protocol:</label>
            <span>{(results as any)?.protocol ?? "-"}</span>
          </div>
        </div>

      </div>
    </div>
  );
}

export default App;