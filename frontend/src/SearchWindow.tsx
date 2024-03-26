import { Button, Input, Label } from "@fluentui/react-components";
import { useCallback, useId, useState } from "react";
import styles from "./SearchWindow.module.css";

export default function SearchWindow({
  onSearch,
}: {
  onSearch?: (query: string) => void;
}) {
  const [query, setQuery] = useState("");
  const onClick = useCallback(() => {
    onSearch && onSearch(query);
  }, [onSearch, query]);

  return (
    <div className={styles.searchContainer}>
      <Input
        value={query}
        onChange={(ev) => setQuery(ev.target.value)}
        placeholder="科目名"
        className={styles.searchInput}
      />
      <Button onClick={onClick} className={styles.searchButton}>
        検索
      </Button>
      <p
        style={{
          gridColumn: "1/7",
          marginBottom: "0px",
        }}
      >
        おみくじ
      </p>
      <Button className={styles.rakutanButton}>らくたん</Button>
      <Button className={styles.onitanButton}>おにたん</Button>
      <Button className={styles.tenButton}>10連</Button>
    </div>
  );
}
