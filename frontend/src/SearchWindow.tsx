import { Button, Input, Label } from "@fluentui/react-components";
import { useId } from "react";
import styles from "./SearchWindow.module.css";

export default function SearchWindow() {
  return (
    <div className={styles.searchContainer}>
      {/* <Label htmlFor={searchId}>科目名</Label> */}
      <div className={styles.inputContainer}>
        <Input placeholder="科目名" />
        <Button>検索</Button>
      </div>
    </div>
  );
}
