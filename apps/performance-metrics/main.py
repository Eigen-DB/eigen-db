import streamlit as st
import pandas as pd

def generate_line_chart(csv_path: str, title: str, desc: str, x_label: str, y_label: str) -> None:
    df = pd.read_csv(csv_path)
    reshaped_df = df.pivot(index="num_vectors", columns="k", values="mean_time_secs")
    reshaped_df.columns = [f"k={col}" for col in reshaped_df.columns]
    st.markdown(f'''
        ### {title}
        {desc}
    ''')
    st.line_chart(
        data=reshaped_df,
        x_label=x_label,
        y_label=y_label,
        use_container_width=False,
        width=700,
        height=500
    )

if __name__ == '__main__':
    st.set_page_config() 
    st.markdown('''
    # EigenDB Performance Metrics
    The data has been collected using a [simple script](https://github.com/Eigen-DB/eigen-db/blob/main/benchmarks/benchmarks.py) for benchmarking EigenDB.            
    ''')

    generate_line_chart(
        csv_path='./data/indexing_mean.csv', 
        title='Mean similarity search time',
        desc='This is the average time to perform similarity search with varying numbers of embeddings and values of k.',
        x_label='Number of embeddings', 
        y_label='Mean time (secs)'
    )    
    generate_line_chart(
        csv_path='./data/inserting_mean.csv',  
        title='Mean embedding insertion time',
        desc='This is the average time to insert an embedding with varying numbers of embeddings and values of k.',
        x_label='Number of embeddings', 
        y_label='Mean time (secs)'
    )
    