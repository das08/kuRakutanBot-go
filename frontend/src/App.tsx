import React from "react";
import logo from "./logo.svg";
import "./App.css";
import { FluentProvider, webLightTheme } from "@fluentui/react-components";
import SearchWindow from "./SearchWindow";
import LectureView from "./LectureView";

function App() {
  return (
    <FluentProvider theme={webLightTheme}>
      <div className="App">
        <h1 className="title">京大落単検索 Web</h1>
        <div>
          <SearchWindow />
        </div>
        <div className="lecture-view">
          <LectureView
            lectureName="講義名"
            faculty="理学部"
            category="専門科目"
            credits={2}
            history={[
              {
                year: 2021,
                rate: 78.0,
                credited: 110,
                total: 141,
              },
              {
                year: 2020,
                rate: 75.5,
                credited: 111,
                total: 147,
              },
            ]}
            rakutanRank="A"
            kakomonLink={undefined}
          />
        </div>
      </div>
    </FluentProvider>
  );
}

export default App;
