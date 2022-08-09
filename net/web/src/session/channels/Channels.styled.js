import styled from 'styled-components';
import Colors from 'constants/Colors';

export const ChannelsWrapper = styled.div`
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  background-color: ${Colors.formBackground};
  min-height: 0;

  .search {
    padding: 8px;

    .filter {
      border: 1px solid ${Colors.divider};
      background-color: ${Colors.white};
      border-radius: 8px;
    }
  }

  .results {
    min-height: 0;
    overflow: scroll;
`;
