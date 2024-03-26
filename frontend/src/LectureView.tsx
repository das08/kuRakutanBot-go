import styles from "./LectureView.module.css";

export default function LectureView({
  lectureName,
  faculty,
  category,
  credits,
  history,
  rakutanRank,
  kakomonLink,
}: {
  lectureName: string;
  faculty?: string;
  category?: string;
  credits?: number;
  history: {
    year: number;
    rate: number;
    credited: number;
    total: number;
  }[];
  rakutanRank: string;
  kakomonLink?: string;
}) {
  return (
    <div className={styles.lectureView}>
      <div className={styles.lectureName}>{lectureName}</div>
      <div className={styles.lectureInfoContainer}>
        <div>開講部局</div>
        <div>{faculty}</div>
        <div>群</div>
        <div>{category}</div>
        <div>単位数</div>
        <div>{credits}</div>
      </div>
      <div className={styles.creditRateContainer}>
        <p className={styles.creditRateTitle}>単位取得率</p>
        <div className={styles.creditRateTable}>
          {history.map((h) => {
            return (
              <>
                <p className={styles.creditYear}>{h.year}年度</p>
                <p className={styles.creditRate}>{h.rate} %</p>
                <p className={styles.creditCount}>
                  ({h.credited}/{h.total})
                </p>
              </>
            );
          })}
        </div>
      </div>
      <div className={styles.otherInfoContainer}>
        <p>らくたん判定</p>
        <p>{rakutanRank}</p>
        <p>過去問</p>
        <p>{kakomonLink}</p>
      </div>
      <div className={styles.infoTextContainer}>
        <p>※単位取得率は「各授業の平均学生在籍数」をもとに表示しています</p>
        <p>
          ※らくたん判定は単位取得率をもとに判定しています。詳しくは「判定詳細」と送信してください。
        </p>
      </div>
    </div>
  );
}
